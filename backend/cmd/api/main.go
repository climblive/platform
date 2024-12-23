package main

import (
	"context"
	"log"
	"log/slog"
	"math/rand"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/climblive/platform/backend/internal/authorizer"
	"github.com/climblive/platform/backend/internal/domain"
	"github.com/climblive/platform/backend/internal/events"
	"github.com/climblive/platform/backend/internal/handlers/rest"
	"github.com/climblive/platform/backend/internal/repository"
	"github.com/climblive/platform/backend/internal/scores"
	"github.com/climblive/platform/backend/internal/usecases"
	"github.com/climblive/platform/backend/internal/utils"
	"github.com/lmittmann/tint"
	"github.com/mattn/go-isatty"
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
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	w := os.Stdout

	logger := slog.New(tint.NewHandler(w, nil))

	slog.SetDefault(slog.New(
		tint.NewHandler(w, &tint.Options{
			Level:      slog.LevelDebug,
			TimeFormat: time.Kitchen,
			NoColor:    !isatty.IsTerminal(w.Fd()),
		}),
	))

	slog.SetDefault(logger)

	var barriers []*sync.WaitGroup

	dbPort, _ := strconv.Atoi(os.Getenv("DB_PORT"))

	repo, err := repository.NewDatabase(
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		dbPort,
		os.Getenv("DB_DATABASE"))
	if err != nil {
		if stack := utils.GetErrorStack(err); stack != "" {
			log.Println(stack)
		}

		panic(err)
	}

	authorizer := authorizer.NewAuthorizer(repo)
	eventBroker := events.NewBroker()
	scoreKeeper := scores.NewScoreKeeper(eventBroker, repo)
	scoreEngineStoreHydrator := &scores.StandardEngineStoreHydrator{Repo: repo}
	scoreEngineManager := scores.NewScoreEngineManager(repo, scoreEngineStoreHydrator, eventBroker)

	barriers = append(barriers,
		scoreKeeper.Run(ctx),
		scoreEngineManager.Run(ctx))

	mux := setupMux(repo, authorizer, eventBroker, scoreKeeper)

	httpServer := &http.Server{
		Addr:    "0.0.0.0:8090",
		Handler: mux,
		BaseContext: func(_ net.Listener) context.Context {
			return ctx
		},
	}

	context.AfterFunc(ctx, func() {
		_ = httpServer.Shutdown(context.Background())
	})

	err = httpServer.ListenAndServe()
	switch err {
	case http.ErrServerClosed:
	default:
		if stack := utils.GetErrorStack(err); stack != "" {
			log.Println(stack)
		}

		panic(err)
	}

	for _, barrier := range barriers {
		barrier.Wait()
	}
}

func setupMux(
	repo *repository.Database,
	authorizer *authorizer.Authorizer,
	eventBroker domain.EventBroker,
	scoreKeeper domain.ScoreKeeper,
) *rest.Mux {
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

	mux := rest.NewMux()
	mux.RegisterMiddleware(rest.CORS)
	mux.RegisterMiddleware(authorizer.Middleware)

	mux.HandleFunc("OPTIONS /", HandleCORSPreFlight)

	rest.InstallContenderHandler(mux, &contenderUseCase)
	rest.InstallContestHandler(mux, &contestUseCase)
	rest.InstallCompClassHandler(mux, &compClassUseCase)
	rest.InstallProblemHandler(mux, &problemUseCase)
	rest.InstallTickHandler(mux, &tickUseCase)
	rest.InstallEventHandler(mux, eventBroker, 10*time.Second)

	return mux
}
