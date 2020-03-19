package kached

import "strings"

// ErrMSG Contains an ErrCode and any additional error message.
type ErrMSG struct {
	code  ErrCode
	stack string
}

// Error returns the string error for an ErrCode.
func (e ErrMSG) Error() string {
	return ErrCodeStrings[e.code] + `: ` + e.stack
}

// msg edits the message string for the ErrMSG.
func (e *ErrMSG) msg(s string) {
	e.stack = s
}

// ErrCode is the corresponding code for an error.
type ErrCode int

// Error returns the string error for an ErrCode.
func (c ErrCode) Error() string {
	return ErrCodeStrings[c]
}

// IsErrCode returns true if the error matches the ErrCode, false otherwise.
func IsErrCode(err error, c ErrCode) bool {
	switch {
	case err == nil && c == ErrNoErr:
		return true
	case err == nil:
		return false
	}
	return strings.HasPrefix(err.Error(), c.Error())
}

// ErrCode Constants:
const (
	ErrNoErr ErrCode = iota
	ErrUnableToCache
	ErrUnableToSave
	ErrUnableToCacheOrSave
	ErrNotFoundCache
	ErrNotFoundDB
	ErrNotFoundCacheOrDB
)

// ErrCodeStrings maps Err Codes to Err Messages.
var ErrCodeStrings = [...]string{
	ErrNoErr:               `no error`,
	ErrUnableToCache:       `unable to cache kv pair`,
	ErrUnableToSave:        `unable to save kv pair`,
	ErrUnableToCacheOrSave: `unable to cache or save kv pair`,
	ErrNotFoundCache:       `kv pair not found in cache`,
	ErrNotFoundDB:          `kv pair not found in database`,
	ErrNotFoundCacheOrDB:   `kv pair not found in cache or database`,
}
