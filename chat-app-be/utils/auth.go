package utils

import (
	"net/http"
)

func GetToken(r *http.Request) string {
	return r.Header.Get("x-auth-token")
}
