package middleWare

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v5"
)

const CookieName = "loggin_token"

func CheckLoggin(next echo.HandlerFunc) echo.HandlerFunc {

	return func(ctx *echo.Context) error {
		cookie, err := ctx.Cookie(CookieName)
		if err != nil {
			return ctx.Redirect(http.StatusSeeOther, "/")
		}

		fmt.Println("Cookie: ", cookie)

		return next(ctx)
	}
}
