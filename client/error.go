package client

import "errors"

var (
	KeyNotFoundError = errors.New("key not found")
	NodesDontAnswer  = errors.New("nodes don't answer")
)
