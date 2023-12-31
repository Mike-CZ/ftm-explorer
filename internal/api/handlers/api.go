package handlers

import (
	"ftm-explorer/internal/api/graphql/resolvers"
	"ftm-explorer/internal/api/graphql/schema"
	"ftm-explorer/internal/logger"
	"net/http"

	"github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/relay"
	"github.com/graph-gophers/graphql-transport-ws/graphqlws"
	"github.com/rs/cors"
)

// ApiHandler constructs and return the API HTTP handlers chain for serving GraphQL API calls.
func ApiHandler(corsOrigins []string, resolver *resolvers.RootResolver, log logger.ILogger) http.Handler {
	// Create new CORS handler and attach the logger into it, so we get information on Debug level if needed
	corsHandler := cors.New(corsOptions(corsOrigins))

	// we don't want to write a method for each type field if it could be matched directly
	opts := []graphql.SchemaOpt{graphql.UseFieldResolvers()}

	// create new parsed GraphQL schema
	s := graphql.MustParseSchema(schema.Schema(), resolver, opts...)

	// return the constructed API handler chain
	return &LoggingHandler{
		log:     log,
		handler: corsHandler.Handler(graphqlws.NewHandlerFunc(s, &relay.Handler{Schema: s})),
	}
}

// corsOptions constructs new set of options for the CORS handler based on provided configuration.
func corsOptions(corsOrigins []string) cors.Options {
	return cors.Options{
		AllowedOrigins: corsOrigins,
		AllowedMethods: []string{"HEAD", "GET", "POST"},
		AllowedHeaders: []string{"Origin", "Accept", "Content-Type", "X-Requested-With"},
		MaxAge:         300,
	}
}
