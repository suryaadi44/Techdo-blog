package utils

import (
	"net/http"
	"time"
)

func GetSessionToken(r *http.Request) (string, bool) {
	loggedIn := true
	c := &http.Cookie{}
	if storedCookie, _ := r.Cookie("session_token"); storedCookie != nil {
		c = storedCookie
	}
	if c.Value == "" {
		loggedIn = false
	}

	return c.Value, loggedIn
}

func DeleteSessionCookie(w *http.ResponseWriter) {
	http.SetCookie(*w, &http.Cookie{
		Name:    "session_token",
		MaxAge:  -1,
		Expires: time.Unix(0, 0),
	})
}
