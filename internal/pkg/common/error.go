package common

import "errors"

var ErrNotFound = errors.New("record not found")

var ErrPermissionDenied = errors.New("permission denied")

var ErrUserNotFound = errors.New("user not found")

var ErrUserAlreadyExist = errors.New("user already exists")
