package model

import "errors"

var (
	ErrPartNotFound       = errors.New("part not found")
	ErrRequestParamsEmpty = errors.New("request params is empty")
	ErrEmptyUUID          = errors.New("empty uuid")
)
