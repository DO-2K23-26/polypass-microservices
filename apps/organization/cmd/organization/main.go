package main

import (
    "github.com/DO-2K23-26/polypass-microservices/organization/internal"
    "log"
)

func main() {
    app, err := internal.NewApp()
    if err != nil {
        log.Fatalf("failed to initialize app: %v", err)
    }

    app.Start()
}
