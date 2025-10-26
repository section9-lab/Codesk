package main

import (
	"context"
	"Codesk/backend/service"
)

// App struct
type App struct {
	ctx          context.Context
	greetService *service.GreetService
	timeService  *service.TimeService
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{
		greetService: service.NewGreetService(),
		timeService:  service.NewTimeService(),
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return a.greetService.Greet(name)
}

// GetCurrentTime 获取当前时间
func (a *App) GetCurrentTime() string {
	return a.timeService.GetCurrentTime()
}