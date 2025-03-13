package kmsTypes

import "errors"

var (
	KmsErrNotReady    = errors.New("kms not ready")
	KmsErrMaxErrTimes = errors.New("network error, max retry times")
)
