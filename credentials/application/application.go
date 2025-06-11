package application

type Application interface {
	Ignite() error
	Stop() error
}
