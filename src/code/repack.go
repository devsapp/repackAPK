package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"repack/oss"
	"strings"

	"github.com/rsc/zipmerge/zip"
)

// print error and exit
func perror(msg string, args ...interface{}) {
	log.Printf(msg, args...)
	os.Exit(1)
}

type readerAt interface {
	ReadAt(buf []byte, off int64) (int, error)
}

func copyData(name string, w http.ResponseWriter, r readerAt, offset, size int64) error {
	buf := make([]byte, size)
	n, err := r.ReadAt(buf, offset)
	log.Printf("%s read %d, actual: %d", name, len(buf), n)
	if (err != nil && err != io.EOF) || n != len(buf) {
		return fmt.Errorf("%s read: %v, n: %d", name, err, n)
	}
	n, err = w.Write(buf)
	if err != nil || n != len(buf) {
		return fmt.Errorf("%s resp write: %v, n: %d", name, err, n)
	}

	return nil
}

type sizeWriter struct {
	io.Writer
	size int64
}

func (sw *sizeWriter) Write(buf []byte) (int, error) {
	n, err := sw.Writer.Write(buf)
	if err != nil {
		return n, err
	}
	sw.size += int64(n)
	return n, err
}

func (sw *sizeWriter) Size() int64 {
	return sw.size
}

type resultInfo struct {
	Offset     int64
	FooterSize int64
}

func repackAPK(fcCtx *FCContext) (*os.File, *resultInfo, error) {
	sourceObject, channelID := fcCtx.SourceObject, fcCtx.ChannelID
	footerFile := fmt.Sprintf("/%s/%s.%s.footer", WORK_DIR_BASE, strings.Replace(sourceObject, "/", "_", -1), channelID)
	resultFile := fmt.Sprintf("/%s/%s.%s.meta", WORK_DIR_BASE, strings.Replace(sourceObject, "/", "_", -1), channelID)

	// try read result file
	buf, err := ioutil.ReadFile(resultFile)
	if err == nil {
		var res resultInfo
		if err := json.Unmarshal(buf, &res); err != nil {
			return nil, nil, err
		}

		file, err := os.Open(footerFile)
		if err != nil {
			return nil, nil, err
		}

		return file, &res, nil
	}
	f, err := os.Create(footerFile)
	if err != nil {
		return nil, nil, err
	}

	offset, size := doRepackAPK(f, fcCtx)
	res := resultInfo{
		Offset:     offset,
		FooterSize: size,
	}
	buf, _ = json.Marshal(res)
	if err := ioutil.WriteFile(resultFile, buf, 0644); err != nil {
		return nil, nil, err
	}

	return f, &res, nil
}

func doRepackAPK(w io.Writer, fcCtx *FCContext) (int64, int64) {
	ossReader, err := oss.NewReader(
		oss.OSSConfig{
			Endpoint:        fcCtx.OSSEndpoint,
			AccessKeyID:     fcCtx.Credentials.AccessKeyID,
			AccessKeySecret: fcCtx.Credentials.AccessKeySecret,
			SecurityToken:   fcCtx.Credentials.SecurityToken,
		}, fcCtx.SourceObject)
	if err != nil {
		perror("oss reader: %v", err)
	}
	objectSize, err := ossReader.Size()
	if err != nil {
		perror("object size: %v", err)
	}

	zipReader, err := zip.NewReader(ossReader, objectSize)
	if err != nil {
		perror("zip reader: %v", err)
	}
	appendOffset := zipReader.AppendOffset()
	log.Printf("append offset: %d", appendOffset)

	err = changeManifest(zipReader, fcCtx)
	if err != nil {
		perror("change manifest: %v", err)
	}
	sizeWriter := &sizeWriter{Writer: w}

	writer := zipReader.Append(sizeWriter)

	// copy cpid file
	if err := copyCPID(writer, fcCtx.ChannelID); err != nil {
		perror("copy cpid: %v", err)
	}
	// copy meta files: MANIFEST.MF/CERT.SF/CERT.RSA
	if err := copyMeta(writer, fcCtx); err != nil {
		perror("copy meta: %v", err)
	}

	writer.Close()

	log.Printf("append offset: %d, footer size: %d", appendOffset, sizeWriter.Size())
	return appendOffset, sizeWriter.Size()
}
