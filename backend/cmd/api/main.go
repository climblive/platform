package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"slices"

	"github.com/climblive/platform/backend/internal/authorizer"
	"github.com/climblive/platform/backend/internal/domain"
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

func HandleCORSPreFlight(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
	w.WriteHeader(http.StatusOK)
}

func main() {
	fmt.Println("Hello, Climbers!")
	ctx := context.Background()

	repo, err := repository.NewDatabase("climblive", "secretpassword", "localhost", "climblive")
	if err != nil {
		if stack := utils.GetErrorStack(err); stack != "" {
			log.Println(stack)
		}
		panic(err)
	}

	authorizer := authorizer.NewAuthorizer()
	eventBroker := events.NewBroker()
	scoreKeeper := scores.NewScoreKeeper(eventBroker)

	go scoreKeeper.Run(ctx)

	engine := scores.NewScoreEngine(1, eventBroker, &scores.HardestProblems{Number: 5}, scores.NewBasicRanker(3))

	go engine.Run(context.Background(), 1)

	problems, err := repo.GetProblemsByContest(ctx, nil, 1)
	if err != nil {
		panic(err)
	}

	for problem := range slices.Values(problems) {
		eventBroker.Dispatch(1, domain.ProblemAddedEvent{
			ProblemID:  problem.ID,
			PointsTop:  problem.PointsTop,
			PointsZone: problem.PointsZone,
			FlashBonus: problem.FlashBonus,
		})
	}

	contenders, err := repo.GetContendersByContest(ctx, nil, 1)
	if err != nil {
		panic(err)
	}

	for contender := range slices.Values(contenders) {
		eventBroker.Dispatch(1, domain.ContenderEnteredEvent{
			ContenderID: contender.ID,
			CompClassID: contender.CompClassID,
		})

		ticks, err := repo.GetTicksByContender(ctx, nil, contender.ID)
		if err != nil {
			panic(err)
		}

		for tick := range slices.Values(ticks) {
			eventBroker.Dispatch(1, domain.AscentRegisteredEvent{
				ContenderID:  contender.ID,
				ProblemID:    tick.ProblemID,
				Top:          tick.Top,
				AttemptsTop:  tick.AttemptsTop,
				Zone:         tick.Zone,
				AttemptsZone: tick.AttemptsTop,
			})
		}
	}

	contenderUseCase := usecases.ContenderUseCase{
		Repo:                      repo,
		Authorizer:                authorizer,
		EventBroker:               eventBroker,
		ScoreKeeper:               scoreKeeper,
		RegistrationCodeGenerator: &registrationCodeGenerator{},
	}

	contestUseCase := usecases.ContestUseCase{
		Repo:        repo,
		ScoreKeeper: scoreKeeper,
	}

	compClassUseCase := usecases.CompClassUseCase{
		Repo: repo,
	}

	problemUseCase := usecases.ProblemUseCase{
		Repo: repo,
	}

	tickUseCase := usecases.TickUseCase{
		Repo:        repo,
		Authorizer:  authorizer,
		EventBroker: eventBroker,
	}

	http.HandleFunc("OPTIONS /", HandleCORSPreFlight)

	rest.InstallContenderHandler(&contenderUseCase)
	rest.InstallContestHandler(&contestUseCase)
	rest.InstallCompClassHandler(&compClassUseCase)
	rest.InstallProblemHandler(&problemUseCase)
	rest.InstallTickHandler(&tickUseCase)
	rest.InstallEventHandler(eventBroker)

	err = http.ListenAndServe("localhost:8090", nil)
	if err != nil {
		if stack := utils.GetErrorStack(err); stack != "" {
			log.Println(stack)
		}
		panic(err)
	}
}
