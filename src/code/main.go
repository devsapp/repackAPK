package main

import (
	"io"
	"log"
	"net/http"
	"os"
	fcoss "repack/oss"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

func repackLocal() {
	CertPEM_PATH = "target/cert/test-cert.pem"
	PrivateKeyPEM_PATH = "target/cert/test-priv.pem"
	WORK_DIR_BASE = "/tmp"

	fcCtx := &FCContext{
		SourceObject: os.Getenv("SOURCE_OBJECT"),
		ChannelID:    os.Getenv("CHANNEL_ID"),
		OSSEndpoint:  os.Getenv("OSS_ENDPOINT"),
		Credentials: Credentials{
			AccessKeyID:     os.Getenv("ACCESS_KEY_ID"),
			AccessKeySecret: os.Getenv("ACCESS_KEY_SECRET"),
		},
		WorkDir: "/tmp",
	}

	f, res, err := repackAPK(fcCtx)
	if err != nil {
		log.Printf("repack error: %v", err)
		return
	}
	defer f.Close()
	log.Printf("res: %+v", res)

	ossReader, err := fcoss.NewReader(
		fcoss.OSSConfig{
			Endpoint:        fcCtx.OSSEndpoint,
			AccessKeyID:     fcCtx.Credentials.AccessKeyID,
			AccessKeySecret: fcCtx.Credentials.AccessKeySecret,
			SecurityToken:   fcCtx.Credentials.SecurityToken,
		}, fcCtx.SourceObject)
	if err != nil {
		log.Printf("read oss: %v", err)
		return
	}
	resp, err := ossReader.Client.GetObject(
		ossReader.Object, oss.Range(0, res.Offset-1))
	if err != nil {
		log.Printf("get object: %v", err)
		return
	}
	defer resp.Close()

	resFile, err := os.Create("/tmp/res.apk")
	if err != nil {
		log.Printf("get object: %v", err)
		return
	}
	defer resFile.Close()
	n, err := io.Copy(resFile, resp)
	if err != nil {
		log.Printf("copy: %v", err)
		return
	}
	log.Printf("copied %d bytes", n)

	if _, err := f.Seek(0, 0); err != nil {
		log.Printf("seek error: %v", err)
		return
	}
	n, err = io.Copy(resFile, f)
	if err != nil {
		log.Printf("copy: %v", err)
		return
	}
	log.Printf("copied %d bytes", n)
}

func main() {
	if os.Getenv("RUN_LOCAL") == "true" {
		repackLocal()
		return
	}

	http.HandleFunc("/", handler)
	http.ListenAndServe(":80", nil)
}
