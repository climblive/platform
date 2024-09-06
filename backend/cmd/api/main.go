package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/climblive/platform/backend/internal/authorizer"
	"github.com/climblive/platform/backend/internal/events"
	"github.com/climblive/platform/backend/internal/handlers/rest"
	"github.com/climblive/platform/backend/internal/repository"
	"github.com/climblive/platform/backend/internal/scores"
	"github.com/climblive/platform/backend/internal/usecases"
	"github.com/climblive/platform/backend/internal/utils"
)

func main() {
	fmt.Println("Hello, Climbers!")

	repo, err := repository.NewDatabase("climblive", "climblive", "localhost", "climblive")
	if err != nil {
		if stack := utils.GetErrorStack(err); stack != "" {
			log.Println(stack)
		}
		panic(err)
	}

	authorizer := authorizer.NewAuthorizer()
	eventBroker := events.NewBroker()
	scoreKeeper := scores.NewScoreKeeper()

	contenderUseCase := usecases.ContenderUseCase{
		Repo:        repo,
		Authorizer:  authorizer,
		EventBroker: eventBroker,
		ScoreKeeper: scoreKeeper,
	}

	rest.InstallContenderHandler(&contenderUseCase)

	err = http.ListenAndServe("localhost:80", nil)
	if err != nil {
		if stack := utils.GetErrorStack(err); stack != "" {
			log.Println(stack)
		}
		panic(err)
	}
}
