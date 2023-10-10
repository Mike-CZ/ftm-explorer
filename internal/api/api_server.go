package api

import (
	"ftm-explorer/internal/api/graphql/resolvers"
	"ftm-explorer/internal/api/handlers"
	"ftm-explorer/internal/api/middlewares"
	"ftm-explorer/internal/config"
	"ftm-explorer/internal/faucet"
	"ftm-explorer/internal/logger"
	"ftm-explorer/internal/repository"
	"net/http"
	"time"
)

// ApiServer represents a GraphQL API server.
// It is responsible for handling GraphQL requests.
type ApiServer struct {
	cfg      *config.ApiServer
	log      logger.ILogger
	srv      *http.Server
	resolver *resolvers.RootResolver
}

// NewApiServer creates a new GraphQL API server.
func NewApiServer(cfg *config.Config, repo repository.IRepository, faucet faucet.IFaucet, log logger.ILogger) *ApiServer {
	apiLogger := log.ModuleLogger("api")
	server := &ApiServer{
		resolver: resolvers.NewResolver(repo, apiLogger, faucet, cfg.Explorer.IsPersisted),
		cfg:      &cfg.Api,
		log:      apiLogger.ModuleLogger("api"),
	}
	server.makeHttpServer()
	return server
}

// Start starts the GraphQL API server.
// It blocks until the server is stopped.
func (api *ApiServer) Start() {
	api.log.Notice("starting API server")
	if err := api.srv.ListenAndServe(); err != nil {
		api.log.Fatalf("failed to start API server: %s", err.Error())
	}
}

// makeHttpServer creates and configures the HTTP server to be used to serve incoming requests
func (api *ApiServer) makeHttpServer() {
	// create request MUXer
	srvMux := http.NewServeMux()

	h := http.TimeoutHandler(
		middlewares.AuthMiddleware(
			handlers.ApiHandler(api.cfg.CorsOrigin, api.resolver, api.log),
		),
		time.Second*time.Duration(api.cfg.ResolverTimeout),
		"Service timeout.",
	)

	srvMux.Handle("/", h)
	srvMux.Handle("/api", h)
	srvMux.Handle("/graphql", h)

	// handle GraphiQL interface
	srvMux.Handle("/graphi", handlers.GraphiHandler(api.cfg.DomainAddress, api.log))

	// create HTTP server to handle our requests
	api.srv = &http.Server{
		Addr:              api.cfg.BindAddress,
		ReadTimeout:       time.Second * time.Duration(api.cfg.ReadTimeout),
		WriteTimeout:      time.Second * time.Duration(api.cfg.WriteTimeout),
		IdleTimeout:       time.Second * time.Duration(api.cfg.IdleTimeout),
		ReadHeaderTimeout: time.Second * time.Duration(api.cfg.HeaderTimeout),
		Handler:           srvMux,
	}
}
