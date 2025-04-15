package application

type Application interface {
	// script to be launched when repository is registered (e.g. for database migration)
	Setup() error

	// script to be launched when repository is deregistered (most of the time for graceful shutdown)
	Shutdown() error

	// script to be launched when application is registered (e.g. for HTTP API)
	Ignite() error

	// script to be launched when application is deregistered (most of the time for graceful shutdown)
	Stop() error
}
