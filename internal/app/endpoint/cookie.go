package endpoint

import (
	"net/http"
	"time"
)

const CookieLiveTime = 20

func (e *EndPoint) setCookie(login string) *http.Cookie {
	cookie := &http.Cookie{
		Name:     "loggin_token",
		Value:    e.s.MakeUUID(),
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Now().Add(CookieLiveTime * time.Minute),
	}

	return cookie
}
