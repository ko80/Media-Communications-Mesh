package registry

import "errors"

const (
	FlagAddStatus uint32 = 1 << iota
	FlagAddConfig
)

var (
	ErrResourceNotFound = errors.New("resource not found")
	ErrTypeCastFailed   = errors.New("type cast failed")
)
