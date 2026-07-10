package endpoint

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
	"net/http"
	"strings"
	"time"
)

type Service interface {
	RegisterUser(login string, password string) error
	LoginUser(login string, password string) error
	GetIdFromDB(login string) (int, error)
	RedisSetKeyValue(key string, value int) error
}

type EndPoint struct {
	s Service
}

func New(svc Service) *EndPoint {
	return &EndPoint{
		s: svc,
	}
}

func (e *EndPoint) LoadMainHTML(ctx *echo.Context) error {
	err := ctx.File("website/LocalCloudMain.html")
	if err != nil {
		fmt.Println(err)
	}

	return err
}

func (e *EndPoint) Register(ctx *echo.Context) error {
	login := ctx.FormValue("login")
	password := ctx.FormValue("password")

	fmt.Println("User try to register: ", login)
	if len(password) > 14 || len(password) < 4 {
		fmt.Println("Password size more than 12 length")
		return ctx.Redirect(http.StatusSeeOther, "/")
	}

	err := e.s.RegisterUser(login, password)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE") {
			return ctx.Redirect(http.StatusSeeOther, "/")
		}
		return errors.New("error with register user, try again later")
	}
	return ctx.Redirect(http.StatusMovedPermanently, "/login.html")
}

func (e *EndPoint) Login(ctx *echo.Context) error {
	login := ctx.FormValue("login")
	password := ctx.FormValue("password")

	fmt.Println("User try to login: ", login)
	fmt.Println("Password for Login: ", password)

	fmt.Println("befora login")
	err := e.s.LoginUser(login, password)
	fmt.Println("after login")
	if err != nil {
		if strings.Contains(err.Error(), "user not found") {

			fmt.Println("Error, account alreade exist: ", err.Error())
			return ctx.Redirect(http.StatusSeeOther, "/")
		}
		fmt.Println("Error with login: ", err)
		return ctx.Redirect(http.StatusSeeOther, "/")
	}

	fmt.Println("adding cookie")

	cookie := new(http.Cookie)
	cookie.Name = "loggin_token"
	cookie.Value = uuid.New().String()
	cookie.Path = "/"

	id, err := e.s.GetIdFromDB(login)
	if err != nil {
		return ctx.Redirect(http.StatusSeeOther, "/")
	}
	err = e.s.RedisSetKeyValue(cookie.Value, id)
	if err != nil {
		return ctx.Redirect(http.StatusSeeOther, "/")
	}
	cookie.Expires = time.Now().Add(20 * time.Minute)
	cookie.HttpOnly = true

	ctx.SetCookie(cookie)

	return ctx.Redirect(http.StatusSeeOther, "/dashboard.html")
}
