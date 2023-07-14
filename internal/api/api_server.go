package api

import (
	"ftm-explorer/internal/api/graphql/resolvers"
	"ftm-explorer/internal/api/handlers"
	"ftm-explorer/internal/config"
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
func NewApiServer(cfg *config.ApiServer, repo repository.IRepository, log logger.ILogger) *ApiServer {
	apiLogger := log.ModuleLogger("api")
	server := &ApiServer{
		resolver: resolvers.NewResolver(repo, apiLogger),
		cfg:      cfg,
		log:      apiLogger,
	}
	server.makeHttpServer()
	return server
}

// Run starts the GraphQL API server.
// It blocks until the server is stopped.
func (api *ApiServer) Run() {
	if err := api.srv.ListenAndServe(); err != nil {
		api.log.Fatalf("Failed to start API server: %s", err.Error())
	}
}

// makeHttpServer creates and configures the HTTP server to be used to serve incoming requests
func (api *ApiServer) makeHttpServer() {
	// create request MUXer
	srvMux := http.NewServeMux()

	h := http.TimeoutHandler(
		handlers.ApiHandler(api.cfg.CorsOrigin, api.resolver, api.log),
		time.Second*time.Duration(api.cfg.ResolverTimeout),
		"Service timeout.",
	)

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
