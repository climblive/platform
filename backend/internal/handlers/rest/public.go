package rest

import "time"

type CreateCompetitorsArguments struct {
	Number int `json:"number"`
}

type StartScoreEngineArguments struct {
	TerminatedBy time.Time `json:"terminatedBy"`
}
