package controllers

import (
	"html/template"
	"log"
	"net/http"
	"os"
)

func Index(res http.ResponseWriter, req *http.Request) {
	tpl, err := template.ParseGlob("templates/html/*.html")
	if err != nil {
		log.Fatalln(err)
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = req.Cookie(os.Getenv("COOKIE_SID"))
	if err == http.ErrNoCookie {
		tpl.ExecuteTemplate(res, "index.html", nil)
		return
	}
	http.Redirect(res, req, "/home", http.StatusSeeOther)
	// name := c.Value
	// log.Println("Found name:", name)
	// tpl.ExecuteTemplate(res, "index.html", name)
}
