package controllers

import (
	"fmt"
	"net/http"
)

func Index(res http.ResponseWriter, req *http.Request) {
	cookie, err := req.Cookie("app_UserSession")
	if err != http.ErrNoCookie {
		fmt.Fprintf(res, "Welcome %s\n", cookie.Value)
		return
	}

	fmt.Fprint(res, "Welcome!\n")
}
