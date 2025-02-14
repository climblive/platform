package rest

import "time"

type CreateContendersArguments struct {
	Number int `json:"number"`
}

type StartScoreEngineArguments struct {
	TerminatedBy time.Time `json:"terminatedBy"`
}
