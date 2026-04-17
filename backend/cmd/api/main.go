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

const appCSP = "default-src 'self'; connect-src 'self' clmb.auth.eu-west-1.amazoncognito.com *.fontawesome.com *.sentry.io data:; style-src 'self' https://fonts.googleapis.com 'unsafe-inline'; font-src 'self' https://fonts.gstatic.com; object-src 'none'; frame-ancestors 'none'; form-action 'none'; base-uri 'self'; img-src 'self' data:; report-uri https://o4509937603641344.ingest.de.sentry.io/api/4509937616093264/security/?sentry_key=019099d850441f60cea5d465e217f768"

const wwwCSP = "default-src 'self'; script-src 'self' 'sha256-jIhoHP5AYEa/rjrf399lCKS/+7hIAc+G1cKDLBSPd7o='; style-src 'self' https://fonts.googleapis.com 'unsafe-inline'; font-src 'self' https://fonts.gstatic.com; frame-ancestors 'none'; form-action 'none'; base-uri 'self'"

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

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	w := os.Stdout

	logger := slog.New(tint.NewHandler(w, nil))

	slog.SetDefault(slog.New(
		tint.NewHandler(w, &tint.Options{
			Level:       slog.LevelDebug,
			TimeFormat:  time.Kitchen,
			NoColor:     !isatty.IsTerminal(w.Fd()),
			AddSource:   false,
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
	appMux.Handle("/api/", http.StripPrefix("/api", noCacheHandler(apiMux)))
	installAppStaticHandlers(appMux)

	wwwMux := http.NewServeMux()
	installWWWStaticHandlers(wwwMux)

	wwwHost := os.Getenv("WWW_HOST")

	tlsConfig := loadTLSConfig()

	handler := maxBytesHandler(newHostHandler(appMux, wwwMux, wwwHost), 1<<20)

	httpServer := &http.Server{
		Addr:                         "0.0.0.0:443",
		Handler:                      handler,
		DisableGeneralOptionsHandler: false,
		TLSConfig:                    tlsConfig,
		ReadTimeout:                  0,
		ReadHeaderTimeout:            0,
		WriteTimeout:                 0,
		IdleTimeout:                  0,
		MaxHeaderBytes:               0,
		TLSNextProto:                 nil,
		ConnState:                    nil,
		ErrorLog:                     nil,
		BaseContext: func(_ net.Listener) context.Context {
			return ctx
		},
		ConnContext: nil,
		HTTP2:       nil,
		Protocols:   nil,
	}

	context.AfterFunc(ctx, func() {
		_ = httpServer.Shutdown(context.Background())
	})

	err = httpServer.ListenAndServeTLS("", "")

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

func loadTLSConfig() *tls.Config {
	type certPair struct {
		cert string
		key  string
	}

	pairs := []certPair{
		{cert: os.Getenv("TLS_APP_CERT_FILE"), key: os.Getenv("TLS_APP_KEY_FILE")},
		{cert: os.Getenv("TLS_WWW_CERT_FILE"), key: os.Getenv("TLS_WWW_KEY_FILE")},
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
		panic("no TLS certificates configured; set TLS_APP_CERT_FILE/TLS_APP_KEY_FILE and/or TLS_WWW_CERT_FILE/TLS_WWW_KEY_FILE")
	}

	return &tls.Config{
		Certificates: certificates,
		MinVersion:   tls.VersionTLS12,
	}
}

type hostHandler struct {
	wwwHost    string
	appHandler http.Handler
	wwwHandler http.Handler
}

func newHostHandler(appHandler, wwwHandler http.Handler, wwwHost string) *hostHandler {
	return &hostHandler{
		wwwHost:    wwwHost,
		appHandler: appHandler,
		wwwHandler: wwwHandler,
	}
}

func (h *hostHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	host := r.Host
	if colonIdx := strings.LastIndex(host, ":"); colonIdx != -1 {
		host = host[:colonIdx]
	}

	if h.wwwHost != "" && host == h.wwwHost {
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
			panic(err)
		}

		rest.InstallStaticHandler(mux, app.basePath, subFS, appCSP)
	}
}

func installWWWStaticHandlers(mux *http.ServeMux) {
	subFS, err := fs.Sub(webAssets, "web/www")
	if err != nil {
		panic(err)
	}

	rest.InstallStaticHandler(mux, "/", subFS, wwwCSP)
}

func noCacheHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		next.ServeHTTP(w, r)
	})
}

func maxBytesHandler(next http.Handler, maxBytes int64) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, maxBytes)
		next.ServeHTTP(w, r)
	})
}
