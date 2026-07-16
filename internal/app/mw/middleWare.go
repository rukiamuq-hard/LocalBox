package middleWare

import (
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

func (mw *MiddleWare) CheckLogin(next echo.HandlerFunc) echo.HandlerFunc {

	return func(ctx *echo.Context) error {
		cookie, err := ctx.Cookie(CookieName)
		if err != nil {
			return ctx.Redirect(http.StatusSeeOther, "/")
		}

		id, err := mw.service.RedisGetValue(cookie.Value)
		if err != nil {
			return ctx.Redirect(http.StatusSeeOther, "/")
		}

		ctx.Set("userID", id)
		return next(ctx)
	}
}
