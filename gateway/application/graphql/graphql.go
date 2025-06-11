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
	"github.com/gofiber/fiber/v2"
	"github.com/optique-dev/modules/graphql/graph"

	"github.com/valyala/fasthttp/fasthttpadaptor"
)

type GraphQL interface {
	Register(app *fiber.App)
	Query() (*fiber.Ctx, error)
	Playground() (*fiber.Ctx, error)
}

type graphqlController struct{}

func wrapHandler(f func(nethttp.ResponseWriter, *nethttp.Request)) func(*fiber.Ctx) {
	return func(ctx *fiber.Ctx) {
		fasthttpadaptor.NewFastHTTPHandler(nethttp.HandlerFunc(f))(ctx.Context())
	}
}

func NewGraphQL() *graphqlController {
	return &graphqlController{}
}

func (g *graphqlController) Query() fiber.Handler {
	h := handler.New(graph.NewExecutableSchema(graph.Config{
		Resolvers: &graph.Resolver{},
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

	app.Get("/playground", g.Playground())
	app.Get("/playground", g.Playground())
}
