package app

import (
	"Umbrella/internal/app/endpoint"          // EndPoint
	middleWare "Umbrella/internal/app/mw"     // MiddleWare
	"Umbrella/internal/app/repository/SQLite" // SQLITE
	"Umbrella/internal/app/service"           // Service
	"fmt"

	"github.com/labstack/echo/v5"
)

type App struct {
	ePoint *endpoint.EndPoint
	servc  *service.Service
	echo   *echo.Echo
	db     *dataBase.DataBase
}

func New() (*App, error) {
	a := &App{}

	a.db = dataBase.New()
	err := a.db.StartDB()
	if err != nil {
		return nil, err
	}

	a.servc = service.New(a.db)
	a.ePoint = endpoint.New(a.servc)
	a.echo = echo.New()
	a.echo.Static("/", "website")

	a.echo.GET("/", a.ePoint.LoadMainHTML)
	a.echo.POST("/register.html", a.ePoint.Register)
	a.echo.POST("/login.html", a.ePoint.Login)
	//upper is can be watching without login

	//down is with cookie secure
	a.echo.GET("/dashboard.html", func(ctx *echo.Context) error { return nil }, middleWare.CheckLoggin)

	return a, nil
}

func (app *App) Run() error {
	fmt.Println("Server running")
	if err := app.echo.Start(":8080"); err != nil {
		return fmt.Errorf("Error to start: %w", err)
	}

	return nil
}

func (app *App) Close() {
	if app.db != nil {
		app.db.CloseDB()
	}

}
