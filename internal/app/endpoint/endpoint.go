package endpoint

import (
	"errors"
	"fmt"
	"github.com/labstack/echo/v5"
	"net/http"
	"strings"
)

type Service interface {
	RegisterUser(login string, password string) error
}

type EndPoint struct {
	s Service
}

func New(svc Service) *EndPoint {
	return &EndPoint{
		s: svc,
	}
}

func (e *EndPoint) LoadLoginReg(ctx *echo.Context) error {
	err := ctx.File("website/login.html")
	if err != nil {
		fmt.Println(err)
	}

	return err
}

func (e *EndPoint) Register(ctx *echo.Context) error {
	login := ctx.FormValue("login")
	password := ctx.FormValue("password")
	if login == "" || password == "" {
		return errors.New("Login/Password is empty")
	}
	fmt.Println("User try to login: ", login)

	err := e.s.RegisterUser(login, password)
	if err != nil {
		fmt.Println(err)
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return ctx.String(http.StatusUnauthorized, "Account already exists")
		}
		return err
	}
	return ctx.Redirect(http.StatusMovedPermanently, "/register")
}
