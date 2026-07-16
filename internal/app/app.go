package app

import (
	"Umbrella/internal/app/endpoint"             // EndPoint
	middleWare "Umbrella/internal/app/mw"        // MiddleWare
	"Umbrella/internal/app/repository/database"  // SQLITE
	rds "Umbrella/internal/app/repository/redis" // Redis
	"Umbrella/internal/app/service"              // Service
	"fmt"

	"github.com/labstack/echo/v5"
)

type App struct {
	ePoint  *endpoint.EndPoint
	servc   *service.Service
	echo    *echo.Echo
	db      *dataBase.DataBase
	rdb     *rds.RedisDB
	midleWR *middleWare.MiddleWare
}

func New() (*App, error) {
	a := &App{}

	a.db = dataBase.New() // SQLITE START
	err := a.db.StartDB()
	if err != nil {
		return nil, err
	}

	a.rdb = rds.New() // REDIS START
	err = a.rdb.StartRedis()
	if err != nil {
		return nil, err
	}

	a.servc = service.New(a.db, a.rdb)
	a.ePoint = endpoint.New(a.servc)
	a.midleWR = middleWare.New(a.servc)
	a.echo = echo.New()
	a.echo.Static("/", "website")

	a.echo.GET("/", a.ePoint.LoadMainHTML)
	a.echo.POST("/register.html", a.ePoint.Register)
	a.echo.POST("/login.html", a.ePoint.Login)
	//upper is can be watching without login

	//down is with cookie secure
	a.echo.GET("/dashboard.html", a.ePoint.LoadDashboard, a.midleWR.CheckLogin)
	a.echo.POST("/api/files", a.ePoint.UploadFile, a.midleWR.CheckLogin)
	a.echo.GET("/api/files", a.ePoint.GetFiles, a.midleWR.CheckLogin)
	a.echo.GET("/download/:id", a.ePoint.DownloadFile, a.midleWR.CheckLogin)

	return a, nil
}

func (app *App) Run() error {
	fmt.Println("[ECHO]: Initialized")
	if err := app.echo.Start(":8080"); err != nil {
		return fmt.Errorf("Error to start: %w", err)
	}

	return nil
}

func (app *App) Close() {
	if app.db != nil && app.rdb != nil {
		app.db.CloseDB()
		app.rdb.Close()
	}
}
