package main

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

// Credentials ...
type Credentials struct {
	AccessKeyID     string
	AccessKeySecret string
	SecurityToken   string
}

// ServiceMeta ...
type ServiceMeta struct {
	ServiceName string
	LogProject  string
	LogStore    string
	Qualifier   string
	VersionID   string
}

// FunctionMeta ...
type FunctionMeta struct {
	Name                  string
	Handler               string
	Memory                int
	Timeout               int
	Initializer           string
	InitializationTimeout int
}

// FCContext ...
type FCContext struct {
	RequestID   string
	Credentials Credentials
	Function    FunctionMeta
	Service     ServiceMeta
	Region      string
	AccountID   string

	SourceObject   string
	ChannelID      string
	NewApkFileName string
	OSSEndpoint    string
	WorkDir        string
	SigFileName    string
}

// NewFromContext ...
func NewFromContext(req *http.Request) (*FCContext, error) {
	mStr := req.Header.Get(fcFunctionMemory)
	m, err := strconv.Atoi(mStr)
	if err != nil {
		m = -1
	}
	tStr := req.Header.Get(fcFunctionTimeout)
	t, err := strconv.Atoi(tStr)
	if err != nil {
		t = -1
	}
	itStr := req.Header.Get(fcInitializationTimeout)
	it, err := strconv.Atoi(itStr)
	if err != nil {
		it = -1
	}
	rid := req.Header.Get(fcRequestID)

	sourceObject := req.URL.Query().Get("src")
	channelID := req.URL.Query().Get("cid")
	ossEndpoint := fmt.Sprintf("http://oss-%s-internal.aliyuncs.com", req.Header.Get(fcRegion))
	bucketAndObject := strings.SplitN(sourceObject, "/", 2)
	if len(bucketAndObject) != 2 {
		return nil, fmt.Errorf("src = %s is invalid, the format is bucket/objectkey", sourceObject)
	}
	_, objectKey := bucketAndObject[0], bucketAndObject[1]
	_, fileName := filepath.Split(objectKey)
	fileSuffix := path.Ext(fileName)
	filenameOnly := strings.TrimSuffix(fileName, fileSuffix)
	newApkFileName := fmt.Sprintf("%s_%s.apk", filenameOnly, channelID)

	workDir := fmt.Sprintf("/%s/%s.%s_workdir", WORK_DIR_BASE, strings.Replace(sourceObject, "/", "_", -1), channelID)
	exist, _ := PathExists(workDir)
	if !exist {
		err := os.MkdirAll(workDir, os.ModePerm)
		if err != nil {
			return nil, fmt.Errorf("fail to create dir: %s; err %v", workDir, err)
		}
	}

	ctx := &FCContext{
		RequestID: rid,
		Credentials: Credentials{
			AccessKeyID:     req.Header.Get(fcAccessKeyID),
			AccessKeySecret: req.Header.Get(fcAccessKeySecret),
			SecurityToken:   req.Header.Get(fcSecurityToken),
		},
		Function: FunctionMeta{
			Name:                  req.Header.Get(fcFunctionName),
			Handler:               req.Header.Get(fcFunctionHandler),
			Memory:                m,
			Timeout:               t,
			Initializer:           req.Header.Get(fcFunctionInitializer),
			InitializationTimeout: it,
		},
		Service: ServiceMeta{
			ServiceName: req.Header.Get(fcServiceName),
			LogProject:  req.Header.Get(fcServiceLogProject),
			LogStore:    req.Header.Get(fcServiceLogstore),
			Qualifier:   req.Header.Get(fcQualifier),
			VersionID:   req.Header.Get(fcVersionID),
		},
		Region:    req.Header.Get(fcRegion),
		AccountID: req.Header.Get(fcAccountID),

		SourceObject:   sourceObject,
		ChannelID:      channelID,
		NewApkFileName: newApkFileName,
		OSSEndpoint:    ossEndpoint,
		WorkDir:        workDir,
		SigFileName:    "",
	}
	return ctx, nil
}
