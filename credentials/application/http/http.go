package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

type Http interface {
	Ignite() error
	Stop() error
	WithHandler(Handler)
}

type http struct {
	listen_addr string
	app         *fiber.App
	handlers    []Handler
}

func NewHttp(config Config) (*http, error) {
	return &http{
		listen_addr: config.ListenAddr,
		app:         fiber.New(),
		handlers:    []Handler{},
	}, nil
}

func (m *http) WithHandler(handler Handler) {
	m.handlers = append(m.handlers, handler)
}

func (m *http) Ignite() error {
	m.app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
	}))

	m.app.Use(logger.New())

	for _, handler := range m.handlers {
		handler.Register(m.app)
	}

	return m.app.Listen(m.listen_addr)
}

func (m *http) Stop() error {
	return m.app.Shutdown()
}
