package core

import (
	"strings"
	"time"
)

// ServerConfiguration struct
type ServerConfiguration struct {
	WhitelistedExtensions []string
	MaximumWidth          int
	LocalBasePath         string
	RemoteBasePath        string
	RemoteBaseURL         string
	DefaultQuality        uint
	UploaderConcurrency   uint
	ProcessorConcurrency  uint
	HTTPTimeout           time.Duration
	Adapters              *Adapters
	Outputs               string
	AWSAccessKeyID        string
	AWSSecretKey          string
	AWSBucket             string
	AWSRegion             string
	MantaURL              string
	MantaUser             string
	MantaKeyID            string
	SDCIdentity           string
	UploaderType          string
	CleanUpTicker         *time.Ticker
	MaxFileAge            time.Duration
	ClearOnlyOriginals    bool
}

func (sc *ServerConfiguration) UploaderIsAws() bool {
	uploader := strings.ToLower(sc.UploaderType)
	if uploader == "aws" || uploader == "s3" {
		return true
	}
	return false
}

func (sc *ServerConfiguration) UploaderIsManta() bool {
	uploader := strings.ToLower(sc.UploaderType)
	if uploader == "manta" {
		return true
	}
	return false
}
