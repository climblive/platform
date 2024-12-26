package domain

import "errors"

var ErrNotFound = errors.New("not found")
var ErrBadState = errors.New("bad state")
var ErrNotAuthorized = errors.New("not authorized")
var ErrNoOwnership = errors.New("no ownership")
var ErrNotAllowed = errors.New("not allowed")
var ErrInsufficientRole = errors.New("insufficient role")
var ErrContestEnded = errors.New("contest ended")
var ErrRepositoryIntegrityViolation = errors.New("repository integrity violation")
var ErrInvalidData = errors.New("invalid data")
var ErrEmptyName = errors.New("name is empty")
var ErrLimitExceeded = errors.New("limit exceeded")
var ErrNotRegistered = errors.New("not registered")
var ErrProblemNotInContest = errors.New("problem not in contest")
