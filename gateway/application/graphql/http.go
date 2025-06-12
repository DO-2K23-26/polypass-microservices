package graphql

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

type HTTP interface {
	Ignite() error
	Stop() error
	// Add more methods here
}

type Handler interface {
	Register(app *fiber.App)
}

type http struct {
	listen_addr string
	app         *fiber.App
	handlers    []Handler
}

func NewHttp(config Config) *http {
	return &http{
		listen_addr: config.ListenAddr,
		app:         fiber.New(),
		handlers:    []Handler{},
	}
}

func (h *http) WithHandler(handler Handler) {
	h.handlers = append(h.handlers, handler)
}

func (h *http) Ignite() error {
	h.app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
	}))
	h.app.Use(logger.New())
	for _, handler := range h.handlers {
		handler.Register(h.app)
	}
	return h.app.Listen(h.listen_addr)
}

func (h *http) Stop() error {
	return h.app.Shutdown()
}
