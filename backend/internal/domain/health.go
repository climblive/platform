package domain

type StatusReporter interface {
	GetStatus() RunnerStatus
}
