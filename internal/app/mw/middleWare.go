package middleWare

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v5"
)

const CookieName = "loggin_token"

type Service interface {
	RedisGetValue(key string) (string, error)
}

type MiddleWare struct {
	service Service
}

func New(svc Service) *MiddleWare {
	return &MiddleWare{service: svc}
}

func (mw *MiddleWare) CheckLoggin(next echo.HandlerFunc) echo.HandlerFunc {

	return func(ctx *echo.Context) error {
		cookie, err := ctx.Cookie(CookieName)
		if err != nil {
			return ctx.Redirect(http.StatusSeeOther, "/")
		}

		fmt.Println("Cookie: ", cookie)

		val, err := mw.service.RedisGetValue(cookie.Value)
		if err != nil {
			return ctx.Redirect(http.StatusSeeOther, "/")
		}
		fmt.Println("id", val)

		return next(ctx)
	}
}
