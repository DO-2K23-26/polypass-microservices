package infrastructure

type Repository interface {
	// script to be launched when repository is registered (e.g. for database migration)
	Setup() error

	// script to be launched when the repository is deregistered (most of the time for graceful shutdown)
	Shutdown() error
}
