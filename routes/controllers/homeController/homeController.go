package homecontroller

import (
	"html/template"
	"log"
	"net/http"
	"os"
	redisauth "webApp/auth/redisAuth"

	"github.com/go-redis/redis/v8"
)

type homeController struct {
	rc *redis.Client
}

func New(rcl *redis.Client) *homeController {
	return &homeController{
		rc: rcl,
	}
}

func (hc *homeController) Home(res http.ResponseWriter, req *http.Request) {
	tpl, err := template.ParseGlob("templates/html/*.html")
	if err != nil {
		log.Fatalln(err)
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	c, err := req.Cookie(os.Getenv("COOKIE_SID"))
	if err == http.ErrNoCookie {
		tpl.ExecuteTemplate(res, "index.html", false)
		return
	}
	user, err := redisauth.GetUserFromSession(c.Value, hc.rc)
	if err != nil {
		log.Println("Error getting user from session:", err)
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Println("Found user:", user)
	tpl.ExecuteTemplate(res, "home.html", user)
}
