package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/DO-2K23-26/polypass-microservices/gateway/application"
	"github.com/DO-2K23-26/polypass-microservices/gateway/infrastructure"
	"github.com/optique-dev/optique"
)

// Cycle is the component in charge of the life cycle of the application
// It is responsible for starting quickly your app and shutting it down gracefully

type Cycle interface {
	Setup() error
	Ignite() error
	Stop() error
}

type cycle struct {
	repos    []infrastructure.Repository
	apps     []application.Application
	shutdown chan os.Signal
}

func NewCycle() *cycle {
	return &cycle{
		shutdown: make(chan os.Signal, 1),
		repos:    []infrastructure.Repository{},
		apps:     []application.Application{},
	}
}

func (c *cycle) AddRepository(repo infrastructure.Repository) {
	c.repos = append(c.repos, repo)
}

func (c *cycle) AddApplications(apps []application.Application) {
	c.apps = append(c.apps, apps...)
}

func (c *cycle) AddRepositories(repos []infrastructure.Repository) {
	c.repos = append(c.repos, repos...)
}

func (c *cycle) AddApplication(app application.Application) {
	c.apps = append(c.apps, app)
}

func (c *cycle) Setup() error {
	if len(c.repos) == 0 {
		optique.Info("No repository to setup")
		return nil
	}
	for _, repository := range c.repos {
		if err := repository.Setup(); err != nil {
			return err
		}
	}
	return nil
}

// Ignite starts the application
func (c *cycle) Ignite() error {
	if len(c.apps) == 0 {
		optique.Info("No application to start")
		return nil
	}

	for _, app := range c.apps {
		go func(app application.Application) {
			err := app.Ignite()
			if err != nil {
				optique.Error(err.Error())
			}
		}(app)
	}

	signal.Notify(c.shutdown, os.Interrupt, syscall.SIGTERM)

	_ = <-c.shutdown

	_ = c.Stop()

	return nil
}

// Stop stops the application
func (c *cycle) Stop() error {
	optique.Info("Stopping applications with graceful shutdown")
	close(c.shutdown)
	for _, app := range c.apps {
		err := app.Stop()
		if err != nil {
			return err
		}
	}
	return nil
}
