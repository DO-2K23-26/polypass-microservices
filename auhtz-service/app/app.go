package app


type App struct {
    // ...
}

func NewApp() *App {
    return &App{}
}


// Perform instanciation to external services/ local services/ repos
func (a *App) Start() error {
    // ...
    return nil 
}


// Perform any data migration the version of the app need
func (a *App) Init() error {
    // ...
    return nil 
}


// Allow to stop the app gracefully
func (a *App) Stop() error {
    // ...
    return nil 
}

