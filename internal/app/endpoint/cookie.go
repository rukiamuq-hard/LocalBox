package endpoint

import (
	"net/http"
	"strconv"
	"time"
)

const CookieLiveTime = 20

func (e *EndPoint) setCookie(login string) (*http.Cookie, *http.Cookie, error) {
	id, err := e.s.GetIdFromDB(login)
	if err != nil {
		return nil, nil, err
	}
	cookie := &http.Cookie{
		Name:     "loggin_token",
		Value:    e.s.MakeUUID(),
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Now().Add(CookieLiveTime * time.Minute),
	}

	cookieID := &http.Cookie{
		Name:     "account-id",
		Value:    strconv.Itoa(id),
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Now().Add(CookieLiveTime * time.Minute),
	}

	return cookie, cookieID, nil
}
