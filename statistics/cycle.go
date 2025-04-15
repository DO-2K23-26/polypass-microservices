package main

import (
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/DO-2K23-26/polypass-microservices/statistics/application"
	"github.com/DO-2K23-26/polypass-microservices/statistics/infrastructure"
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
	wg       sync.WaitGroup
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
		slog.Info("No repository to setup")
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
		slog.Info("No application to start")
		return nil
	}

	// Start all applications
	for _, app := range c.apps {
		c.wg.Add(1)
		go func(app application.Application) {
			defer c.wg.Done()
			err := app.Ignite()
			if err != nil {
				slog.Error(err.Error())
			}
		}(app)
	}

	// Wait for shutdown signal
	signal.Notify(c.shutdown, os.Interrupt, syscall.SIGTERM)
	<-c.shutdown

	// Stop all applications
	return c.Stop()
}

// Stop stops the application
func (c *cycle) Stop() error {
	slog.Info("Stopping applications with graceful shutdown")
	close(c.shutdown)

	// Wait for all applications to stop
	c.wg.Wait()

	for _, app := range c.apps {
		err := app.Stop()
		if err != nil {
			return err
		}
	}
	return nil
}
