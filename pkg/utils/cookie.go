package utils

import "net/http"

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
