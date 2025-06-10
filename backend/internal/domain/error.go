package domain

import "errors"

var ErrNotFound = errors.New("not found")
var ErrDuplicate = errors.New("duplicate")
var ErrBadState = errors.New("bad state")
var ErrNotAuthenticated = errors.New("not authenticated")
var ErrNotAuthorized = errors.New("not authorized")
var ErrNoOwnership = errors.New("no ownership")
var ErrNotAllowed = errors.New("not allowed")
var ErrInsufficientRole = errors.New("insufficient role")
var ErrContestNotStarted = errors.New("contest not started")
var ErrContestEnded = errors.New("contest ended")
var ErrRepositoryIntegrityViolation = errors.New("repository integrity violation")
var ErrInvalidData = errors.New("invalid data")
var ErrEmptyName = errors.New("name is empty")
var ErrLimitExceeded = errors.New("limit exceeded")
var ErrNotRegistered = errors.New("not registered")
var ErrProblemNotInContest = errors.New("problem not in contest")
var ErrAllWinnersDrawn = errors.New("all winners drawn")
