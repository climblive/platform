package domain

import "errors"

var ErrNotFound = errors.New("not found")
var ErrBadState = errors.New("bad state")
var ErrNoOwnership = errors.New("no ownership")
var ErrNotAllowed = errors.New("not allowed")
var ErrInsufficientRole = errors.New("insufficient role")
var ErrContestEnded = errors.New("contest ended")
var ErrRepositoryFailure = errors.New("repository failure")
