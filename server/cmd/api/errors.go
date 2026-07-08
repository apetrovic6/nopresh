package main

import (
	"errors"

	"connectrpc.com/connect"
)

var (
	ConnErrMissingAuthInfo = connect.NewError(connect.CodeUnauthenticated, errors.New("missing auth info"))
)
