package domain

import "errors"

var ErrNotFound = errors.New("not found")
var ErrBadState = errors.New("bad state")
var ErrPermissionDenied = errors.New("permission denied")
var ErrNotAllowed = errors.New("not allowed")
var ErrContestEnded = errors.New("contest ended")
var ErrRepositoryFailure = errors.New("repository failure")
