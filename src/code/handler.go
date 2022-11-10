package main

import (
	"fmt"
	"log"
	"net/http"
	"repack/oss"
	"strconv"
	"strings"
)

func handleError(w http.ResponseWriter, err error) {
	w.WriteHeader(400)
	log.Printf("handle error: %v", err)
	fmt.Fprintf(w, "error: %v", err)
}

func parseRange(r string) (int64, int64, error) {
	parts := strings.Split(r, "=")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("invalid range: %s", r)
	}
	parts = strings.Split(parts[1], "-")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("invalid range: %s", r)
	}
	beginPos, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return 0, 0, err
	}
	endPos, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return 0, 0, err
	}

	if sz := endPos + 1 - beginPos; sz <= 0 || sz > 50*1024*1024 {
		return 0, 0, fmt.Errorf("invalid range: %s", r)
	}

	return beginPos, endPos + 1, nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	fcCtx, err := NewFromContext(r)
	log.Printf("fcContext=%v", fcCtx)
	if err != nil {
		handleError(w, fmt.Errorf("fail to NewFromContext due to  %v", err))
		return
	}
	switch r.Method {
	case "HEAD":
		f, res, err := repackAPK(fcCtx)
		if err != nil {
			handleError(w, err)
			return
		}
		defer f.Close()
		w.Header().Set("Accept-Ranges", "bytes")
		w.Header().Set("Content-Length", fmt.Sprintf("%d", res.Offset+res.FooterSize))
		w.WriteHeader(200)
		return
	case "GET":
		log.Printf("range: %s", r.Header.Get("Range"))
		beginPos, endPos, err := parseRange(r.Header.Get("Range"))
		if err != nil {
			log.Printf("parse range error: %v", err)
			handleError(w, err)
			return
		}
		f, res, err := repackAPK(fcCtx)
		if err != nil {
			handleError(w, err)
			return
		}
		defer f.Close()
		if endPos > res.Offset+res.FooterSize {
			endPos = res.Offset + res.FooterSize
		}
		w.Header().Set("Accept-Ranges", "bytes")
		w.Header().Set("Cache-Control", "max-age=604800") // tell CDN to cache 7 days
		w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, fcCtx.NewApkFileName))
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", beginPos, endPos-1, res.Offset+res.FooterSize))
		w.Header().Set("Content-Length", fmt.Sprintf("%d", endPos-beginPos))
		if endPos-beginPos < res.Offset+res.FooterSize {
			w.WriteHeader(206)
		}
		// need read from oss
		ossBegin, ossEnd := int64(-1), int64(-1)
		if beginPos < res.Offset {
			log.Printf("Range: %s, read oss, beginPos: %d, offset: %d", r.Header.Get("Range"), beginPos, res.Offset)
			ossBegin = beginPos
			ossEnd = endPos
			if res.Offset < ossEnd {
				ossEnd = res.Offset
			}
			ossReader, err := oss.NewReader(
				oss.OSSConfig{
					Endpoint:        fcCtx.OSSEndpoint,
					AccessKeyID:     fcCtx.Credentials.AccessKeyID,
					AccessKeySecret: fcCtx.Credentials.AccessKeySecret,
					SecurityToken:   fcCtx.Credentials.SecurityToken,
				}, fcCtx.SourceObject)
			if err != nil {
				handleError(w, fmt.Errorf("oss reader: %v", err))
				return
			}
			err = copyData("oss", w, ossReader, ossBegin, ossEnd-ossBegin)
			if err != nil {
				log.Printf("copy error: %v", err)
				handleError(w, err)
				return
			}
		}
		if endPos > res.Offset {
			log.Printf("Range: %s, read file, endPos: %d, offset: %d", r.Header.Get("Range"), endPos, res.Offset)
			fileBegin := int64(0)
			if beginPos > res.Offset {
				fileBegin = beginPos - res.Offset
			}
			err := copyData("file", w, f, fileBegin, endPos-res.Offset-fileBegin)
			if err != nil {
				log.Printf("copy error: %v", err)
				handleError(w, err)
				return
			}
		}
		return
	default:
		handleError(w, fmt.Errorf("method %s not supported", r.Method))
	}
}
