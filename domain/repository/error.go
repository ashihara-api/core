package repository

import (
	"errors"
)

var (
	// ErrAreadyExist is error data is found in datasource when requests to create
	ErrAreadyExist = errors.New("data is already exist")
	// ErrNotExist is error data is not found in datasource
	ErrNotExist = errors.New("no data in request set")
	// ErrNoPermission ...
	ErrNoPermission = errors.New("no permission")
)
