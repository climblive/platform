package domain

import "errors"

var ErrNotFound = errors.New("not found")
var ErrSuperseded = errors.New("superseded")
var ErrArchived = errors.New("archived")
var ErrDuplicate = errors.New("duplicate")
var ErrBadState = errors.New("bad state")
var ErrNotAuthenticated = errors.New("not authenticated")
var ErrNotAuthorized = errors.New("not authorized")
var ErrNoOwnership = errors.New("no ownership")
var ErrNotAllowed = errors.New("not allowed")
var ErrInsufficientRole = errors.New("insufficient role")
var ErrCompetitionNotStarted = errors.New("competition not started")
var ErrCompetitionEnded = errors.New("competition ended")
var ErrRepositoryIntegrityViolation = errors.New("repository integrity violation")
var ErrInvalidData = errors.New("invalid data")
var ErrEmptyName = errors.New("name is empty")
var ErrLimitExceeded = errors.New("limit exceeded")
var ErrNotRegistered = errors.New("not registered")
var ErrProblemNotInCompetition = errors.New("problem not in competition")
var ErrAllWinnersDrawn = errors.New("all winners drawn")
var ErrExpired = errors.New("expired")
