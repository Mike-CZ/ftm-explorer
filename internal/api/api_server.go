package api

import (
	"ftm-explorer/internal/api/graphql/resolvers"
	"ftm-explorer/internal/api/handlers"
	"ftm-explorer/internal/api/middlewares"
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

	// create handler for api requests
	h := handlers.ApiHandler(api.cfg.CorsOrigin, api.resolver, api.log)

	// if jwt authorization is enabled, wrap the handler with the middleware
	if api.cfg.Jwt.Enabled {
		h = middlewares.JwtMiddleware(
			h,
			api.log,
			api.cfg.Jwt.Secret,
			api.cfg.Jwt.Version,
		)
	}

	// add timeout handler
	h = http.TimeoutHandler(h, time.Second*time.Duration(api.cfg.ResolverTimeout), "Service timeout.")

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
