package endpoint

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v5"
)

func (e *EndPoint) Register(ctx *echo.Context) error {
	login := ctx.FormValue("login")
	password := ctx.FormValue("password")

	fmt.Println("User try to register: ", login)
	if len(password) > 14 || len(password) < 4 {
		fmt.Println("Password size more than 14 length")
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

	err := e.s.LoginUser(login, password)
	if err != nil {
		if strings.Contains(err.Error(), "user not found") {

			fmt.Println("Error, account alreade exist: ", err.Error())
			return ctx.Redirect(http.StatusSeeOther, "/")
		}
		fmt.Println("Error with login: ", err)
		return ctx.Redirect(http.StatusSeeOther, "/")
	}

	cookie, cookieID, err := e.setCookie(login)
	if err != nil {
		return ctx.Redirect(http.StatusSeeOther, "/")
	}
	id, err := e.s.GetIdFromDB(login)
	if err != nil {
		return ctx.Redirect(http.StatusSeeOther, "/")
	}
	err = e.s.RedisSetKeyValue(cookie.Value, id)
	if err != nil {
		return ctx.Redirect(http.StatusSeeOther, "/")
	}

	ctx.SetCookie(cookie)
	ctx.SetCookie(cookieID)

	return ctx.Redirect(http.StatusSeeOther, "/dashboard.html")
}
