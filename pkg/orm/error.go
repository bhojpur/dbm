package orm

import (
	"errors"
)

var (
	// ErrPtrSliceType represents a type error
	ErrPtrSliceType = errors.New("A point to a slice is needed")
	// ErrParamsType params error
	ErrParamsType = errors.New("Params type error")
	// ErrTableNotFound table not found error
	ErrTableNotFound = errors.New("Table not found")
	// ErrUnSupportedType unsupported error
	ErrUnSupportedType = errors.New("Unsupported type error")
	// ErrNotExist record does not exist error
	ErrNotExist = errors.New("Record does not exist")
	// ErrCacheFailed cache failed error
	ErrCacheFailed = errors.New("Cache failed")
	// ErrConditionType condition type unsupported
	ErrConditionType = errors.New("Unsupported condition type")
)
