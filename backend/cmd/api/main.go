package main

import (
	"context"
	"embed"
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

	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

const defaultScoreEngineMaxLifetime = 24 * time.Hour

type registrationCodeGenerator struct {
}

func (g *registrationCodeGenerator) Generate(length int) string {
	const characters = "ABCDEFGHIJKLMNPQRSTUVWXYZ123456789"
	var code []rune

	for range length {
		code = append(code, []rune(characters)[rand.Intn(len(characters))])
	}

	return string(code)
}

func HandleCORSPreFlight(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE, PATCH")
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

	database, err := repository.NewDatabase(
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

	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("mysql"); err != nil {
		panic(err)
	}

	if err := goose.Up(database.Handle, "migrations"); err != nil {
		panic(err)
	}

	jwtDecoder, err := authorizer.NewStandardJWTDecoder()
	if err != nil {
		if stack := utils.GetErrorStack(err); stack != "" {
			log.Println(stack)
		}

		panic(err)
	}

	authorizer := authorizer.NewAuthorizer(database, jwtDecoder)
	eventBroker := events.NewBroker()
	scoreKeeper := scores.NewScoreKeeper(eventBroker, database)
	scoreEngineStoreHydrator := &scores.StandardEngineStoreHydrator{Repo: database}

	scoreEngineMaxLifetime := getScoreEngineMaxLifetime()
	slog.Info("score engine maximum lifetime cap enabled", "max_lifetime", scoreEngineMaxLifetime)

	scoreEngineManager := scores.NewScoreEngineManager(database, scoreEngineStoreHydrator, eventBroker, scoreEngineMaxLifetime)

	barriers = append(barriers,
		scoreKeeper.Run(ctx, scores.WithPanicRecovery()),
		scoreEngineManager.Run(ctx, scores.WithPanicRecovery()))

	mux := setupMux(database, authorizer, eventBroker, scoreKeeper, &scoreEngineManager)

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

func getScoreEngineMaxLifetime() time.Duration {
	env := "SCORE_ENGINE_MAX_LIFETIME"
	maxLifetime := defaultScoreEngineMaxLifetime

	if value, present := os.LookupEnv(env); present {
		lifetime, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			slog.Warn("discarding non-numeric environment variable", "env", env, "error", err)
		} else {
			maxLifetime = time.Duration(lifetime) * time.Second
		}
	}

	return maxLifetime
}

func setupMux(
	repo *repository.Database,
	authorizer *authorizer.Authorizer,
	eventBroker domain.EventBroker,
	scoreKeeper domain.ScoreKeeper,
	scoreEngineManager *scores.ScoreEngineManager,
) *rest.Mux {
	contenderUseCase := usecases.ContenderUseCase{
		Repo:                      repo,
		Authorizer:                authorizer,
		EventBroker:               eventBroker,
		ScoreKeeper:               scoreKeeper,
		RegistrationCodeGenerator: &registrationCodeGenerator{},
	}

	contestUseCase := usecases.ContestUseCase{
		Authorizer:         authorizer,
		Repo:               repo,
		ScoreKeeper:        scoreKeeper,
		ScoreEngineManager: scoreEngineManager,
	}

	compClassUseCase := usecases.CompClassUseCase{
		Authorizer: authorizer,
		Repo:       repo,
	}

	problemUseCase := usecases.ProblemUseCase{
		Repo:        repo,
		Authorizer:  authorizer,
		EventBroker: eventBroker,
	}

	tickUseCase := usecases.TickUseCase{
		Repo:        repo,
		Authorizer:  authorizer,
		EventBroker: eventBroker,
	}

	scoreEngineUseCase := usecases.ScoreEngineUseCase{
		Repo:               repo,
		Authorizer:         authorizer,
		ScoreEngineManager: scoreEngineManager,
	}

	raffleUseCase := usecases.RaffleUseCase{
		Repo:       repo,
		Authorizer: authorizer,
	}

	userUseCase := usecases.UserUseCase{
		Repo:       repo,
		Authorizer: authorizer,
	}

	organizerUseCase := usecases.OrganizerUseCase{
		Repo:       repo,
		Authorizer: authorizer,
	}

	mux := rest.NewMux()
	mux.RegisterMiddleware(rest.CORS)
	mux.RegisterMiddleware(authorizer.Middleware)

	mux.HandleFunc("OPTIONS /", HandleCORSPreFlight)

	rest.InstallContenderHandler(mux, &contenderUseCase)
	rest.InstallContestHandler(mux, &contestUseCase, &compClassUseCase, &tickUseCase, &problemUseCase)
	rest.InstallCompClassHandler(mux, &compClassUseCase)
	rest.InstallProblemHandler(mux, &problemUseCase)
	rest.InstallTickHandler(mux, &tickUseCase)
	rest.InstallEventHandler(mux, eventBroker, 10*time.Second)
	rest.InstallScoreEngineHandler(mux, &scoreEngineUseCase)
	rest.InstallRaffleHandler(mux, &raffleUseCase)
	rest.InstallUserHandler(mux, &userUseCase)
	rest.InstallOrganizerHandler(mux, &organizerUseCase)

	return mux
}
