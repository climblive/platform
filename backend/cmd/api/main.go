package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"

	"github.com/climblive/platform/backend/internal/authorizer"
	"github.com/climblive/platform/backend/internal/events"
	"github.com/climblive/platform/backend/internal/handlers/rest"
	"github.com/climblive/platform/backend/internal/repository"
	"github.com/climblive/platform/backend/internal/scores"
	"github.com/climblive/platform/backend/internal/usecases"
	"github.com/climblive/platform/backend/internal/utils"
)

type registrationCodeGenerator struct {
}

func (g *registrationCodeGenerator) Generate(length int) string {
	const characters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var code []rune

	for range length {
		code = append(code, []rune(characters)[rand.Intn(len(characters))])
	}

	return string(code)
}

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
		Repo:                      repo,
		Authorizer:                authorizer,
		EventBroker:               eventBroker,
		ScoreKeeper:               scoreKeeper,
		RegistrationCodeGenerator: &registrationCodeGenerator{},
	}

	contestUseCase := usecases.ContestUseCase{
		Repo: repo,
	}

	compClassUseCase := usecases.CompClassUseCase{
		Repo: repo,
	}

	problemUseCase := usecases.ProblemUseCase{
		Repo: repo,
	}

	rest.InstallContenderHandler(&contenderUseCase)
	rest.InstallContestHandler(&contestUseCase)
	rest.InstallCompClassHandler(&compClassUseCase)
	rest.InstallProblemHandler(&problemUseCase)

	err = http.ListenAndServe("localhost:80", nil)
	if err != nil {
		if stack := utils.GetErrorStack(err); stack != "" {
			log.Println(stack)
		}
		panic(err)
	}
}
