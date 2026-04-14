package main

import (
	"context"
	"crypto/tls"
	"embed"
	"io/fs"
	"log"
	"log/slog"
	"math/rand"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
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
	"github.com/google/uuid"
	"github.com/lmittmann/tint"
	"github.com/mattn/go-isatty"

	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

//go:embed all:web
var webAssets embed.FS

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

type uuidGenerator struct {
}

func (g *uuidGenerator) Generate() uuid.UUID {
	return uuid.New()
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
			AddSource:  false,
			ReplaceAttr: nil,
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

	apiMux := setupAPIMux(database, authorizer, eventBroker, scoreKeeper, &scoreEngineManager)

	appMux := http.NewServeMux()
	appMux.Handle("/api/", http.StripPrefix("/api", apiMux))
	appMux.HandleFunc("OPTIONS /api/", HandleCORSPreFlight)
	installAppStaticHandlers(appMux)

	wwwMux := http.NewServeMux()
	installWWWStaticHandlers(wwwMux)

	handler := newHostHandler(appMux, wwwMux)

	tlsConfig, tlsEnabled := loadTLSConfig()

	listenAddr := "0.0.0.0:8090"
	if tlsEnabled {
		listenAddr = "0.0.0.0:443"
	}

	httpServer := &http.Server{
		Addr:      listenAddr,
		Handler:   handler,
		TLSConfig: tlsConfig,
		BaseContext: func(_ net.Listener) context.Context {
			return ctx
		},
	}

	context.AfterFunc(ctx, func() {
		_ = httpServer.Shutdown(context.Background())
	})

	if tlsEnabled {
		slog.Info("TLS enabled, starting HTTPS server", "addr", listenAddr)

		redirectServer := &http.Server{
			Addr: "0.0.0.0:80",
			Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				host := r.Host
				if host == "" {
					http.Error(w, "Bad Request", http.StatusBadRequest)
					return
				}

				target := "https://" + host + r.URL.Path
				if r.URL.RawQuery != "" {
					target += "?" + r.URL.RawQuery
				}
				http.Redirect(w, r, target, http.StatusMovedPermanently)
			}),
		}

		context.AfterFunc(ctx, func() {
			_ = redirectServer.Shutdown(context.Background())
		})

		go func() {
			if err := redirectServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				slog.Error("HTTP redirect server error", "error", err)
			}
		}()

		err = httpServer.ListenAndServeTLS("", "")
	} else {
		slog.Info("TLS not configured, starting HTTP server", "addr", listenAddr)
		err = httpServer.ListenAndServe()
	}

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

func setupAPIMux(
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
		EventBroker:        eventBroker,
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
		Repo:        repo,
		Authorizer:  authorizer,
		EventBroker: eventBroker,
	}

	userUseCase := usecases.UserUseCase{
		Repo:       repo,
		Authorizer: authorizer,
	}

	organizerUseCase := usecases.OrganizerUseCase{
		Repo:          repo,
		Authorizer:    authorizer,
		UUIDGenerator: &uuidGenerator{},
	}

	mux := rest.NewMux()
	mux.RegisterMiddleware(rest.CORS)
	mux.RegisterMiddleware(authorizer.Middleware)

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

func loadTLSConfig() (*tls.Config, bool) {
	type certPair struct {
		cert string
		key  string
	}

	pairs := []certPair{
		{os.Getenv("TLS_APP_CERT_FILE"), os.Getenv("TLS_APP_KEY_FILE")},
		{os.Getenv("TLS_WWW_CERT_FILE"), os.Getenv("TLS_WWW_KEY_FILE")},
	}

	var certificates []tls.Certificate
	for _, p := range pairs {
		if p.cert == "" || p.key == "" {
			continue
		}

		cert, err := tls.LoadX509KeyPair(p.cert, p.key)
		if err != nil {
			slog.Error("failed to load TLS certificate", "cert", p.cert, "key", p.key, "error", err)
			panic(err)
		}

		certificates = append(certificates, cert)
		slog.Info("loaded TLS certificate", "cert", p.cert)
	}

	if len(certificates) == 0 {
		return nil, false
	}

	return &tls.Config{
		Certificates: certificates,
	}, true
}

type hostHandler struct {
	appHandler http.Handler
	wwwHandler http.Handler
}

func newHostHandler(appHandler, wwwHandler http.Handler) *hostHandler {
	return &hostHandler{
		appHandler: appHandler,
		wwwHandler: wwwHandler,
	}
}

func (h *hostHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	host := r.Host
	if colonIdx := strings.LastIndex(host, ":"); colonIdx != -1 {
		host = host[:colonIdx]
	}

	if h.wwwHandler != nil && !strings.HasSuffix(host, ".app") {
		h.wwwHandler.ServeHTTP(w, r)
		return
	}

	h.appHandler.ServeHTTP(w, r)
}

func installAppStaticHandlers(mux *http.ServeMux) {
	apps := []struct {
		basePath string
		subDir   string
	}{
		{"/admin", "web/admin"},
		{"/scoreboard", "web/scoreboard"},
		{"/", "web/scorecard"},
	}

	for _, app := range apps {
		subFS, err := fs.Sub(webAssets, app.subDir)
		if err != nil {
			slog.Debug("skipping static handler, directory not found", "path", app.subDir)
			continue
		}

		rest.InstallStaticHandler(mux, app.basePath, subFS)
	}
}

func installWWWStaticHandlers(mux *http.ServeMux) {
	subFS, err := fs.Sub(webAssets, "web/www")
	if err != nil {
		slog.Debug("skipping www static handler, directory not found")
		return
	}

	rest.InstallStaticHandler(mux, "/", subFS)
}
