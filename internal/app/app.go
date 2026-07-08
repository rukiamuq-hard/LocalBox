package app

import (
	"Umbrella/internal/app/endpoint"
	middleWare "Umbrella/internal/app/mw"
	"Umbrella/internal/app/repository"
	"Umbrella/internal/app/service"
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

	a.echo.Use(middleWare.RoleCheck)
	a.echo.GET("/", a.ePoint.LoadLoginReg)
	a.echo.POST("/register", a.ePoint.Register)

	return a, nil
}

func (app *App) Run() error {
	fmt.Println("Server running")
	if err := app.echo.Start(":8080"); err != nil {
		return fmt.Errorf("Error to start: %w", err)
	}

	return nil
}
