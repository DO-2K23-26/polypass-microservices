package graphql

import (
	"context"
	nethttp "net/http"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/DO-2K23-26/polypass-microservices/gateway/core"
	"github.com/DO-2K23-26/polypass-microservices/gateway/graph"
	"github.com/gofiber/fiber/v2"

	"github.com/valyala/fasthttp/fasthttpadaptor"
)

type GraphQL interface {
	Register(app *fiber.App)
	Query() (*fiber.Ctx, error)
	Playground() (*fiber.Ctx, error)
}

type graphqlController struct {
	organizationService core.OrganizationService
}

// NewGraphQL creates a new GraphQL controller with organization and search services
func NewGraphQL(organizationService core.OrganizationService) *graphqlController {
	return &graphqlController{organizationService: organizationService}
}

func (g *graphqlController) Query() fiber.Handler {
	h := handler.New(graph.NewExecutableSchema(graph.Config{
		Resolvers: &graph.Resolver{OrganizationsService: g.organizationService},
	}))
	h.AddTransport(transport.Websocket{
		KeepAlivePingInterval: 10 * time.Second,
	})
	h.AddTransport(transport.GET{})
	h.AddTransport(transport.POST{})
	h.AddTransport(transport.Options{})

	// Add the introspection middleware.
	h.Use(extension.Introspection{})

	h.AroundFields(func(ctx context.Context, next graphql.Resolver) (res any, err error) {
		res, err = next(ctx)
		return res, err
	})
	return func(ctx *fiber.Ctx) error {
		wrapHandler(h.ServeHTTP)(ctx)
		return nil
	}

}

func (g *graphqlController) Playground() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		wrapHandler(playground.Handler("GraphQL playground", "/query"))(ctx)
		return nil
	}
}

func (g *graphqlController) Register(app *fiber.App) {
	app.Post("/graphql", g.Query())
	app.Get("/graphql", g.Query())
	app.Options("/graphql", g.Query())


	app.Post("/query", g.Query())
	app.Get("/query", g.Query())
	app.Options("/query", g.Query())

	app.Get("/playground", g.Playground())
	app.Get("/playground", g.Playground())
}

// wrapHandler adapts a net/http handler to Fiber
func wrapHandler(f func(nethttp.ResponseWriter, *nethttp.Request)) func(*fiber.Ctx) {
	return func(ctx *fiber.Ctx) {
		fasthttpadaptor.NewFastHTTPHandler(nethttp.HandlerFunc(f))(ctx.Context())
	}
}
