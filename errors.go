package openuri

import "errors"

var (
	ErrNotAbsolute          = errors.New("Provided URL is not absolute")
	ErrProtocolNotSupported = errors.New("requested protocol is not supported")
	ErrLocalInvalidHost     = errors.New("invalid file: hostname provided")
)
