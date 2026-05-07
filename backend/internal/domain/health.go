package domain

type StatusReporter interface {
	GetStatus() ServiceStatus
}
