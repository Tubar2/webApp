package sessioncontroller

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"webApp/auth"

	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
)

type sessionController struct {
	userCollection *mongo.Collection
	redisClient    *redis.Client
}

func New(cl *mongo.Client, rc *redis.Client) *sessionController {
	return &sessionController{
		userCollection: cl.Database("mydb").Collection("users"),
		redisClient:    rc,
	}
}

func (rc *sessionController) Logout(res http.ResponseWriter, req *http.Request) {
	log.Println("Logout handler called")

	c_name := os.Getenv("COOKIE_SID")
	cookie, err := req.Cookie(c_name)
	if err == http.ErrNoCookie {
		log.Println("No cookie during logout")
		http.Error(res, "No session", http.StatusUnprocessableEntity)
		return
	}
	log.Println("Cookie value:", cookie.Value)
	dels, err := auth.Logout(rc.redisClient, cookie.Value)
	if err != nil {
		log.Println("Error during logout:", err)
		http.Error(res, "Error logging out", http.StatusInternalServerError)
		return
	}
	if dels == 0 {
		log.Println("No fields deleted:", err)
		http.Error(res, "Error loggin out: no session", http.StatusInternalServerError)
		return
	}

	cookie = &http.Cookie{
		Name:   cookie.Name,
		MaxAge: -1,
	}
	http.SetCookie(res, cookie)

	http.Redirect(res, req, "/", http.StatusSeeOther)
	res.Write([]byte("Logged out!\n"))
}

func (rc *sessionController) Login(res http.ResponseWriter, req *http.Request) {
	log.Println("Login handler called")

	req.ParseMultipartForm(0)

	uname := req.FormValue("username")
	password := req.FormValue("password")

	sID, err := auth.SignIn(uname, password, rc.userCollection, rc.redisClient)
	if err != nil {
		log.Println("Error signing in:", err)
		http.Error(res, "Error starting session", http.StatusInternalServerError)
		return
	}
	c := &http.Cookie{
		Name:  os.Getenv("COOKIE_SID"),
		Value: sID,
	}
	http.SetCookie(res, c)
	http.Redirect(res, req, "/", http.StatusSeeOther)
}
func (rc *sessionController) LoginPage(res http.ResponseWriter, req *http.Request) {
	tpl, err := template.ParseGlob("templates/html/*.html")
	if err != nil {
		log.Fatalln(err)
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	tpl.ExecuteTemplate(res, "login.html", nil)
}

func (rc *sessionController) RegisterPage(res http.ResponseWriter, req *http.Request) {
	tpl, err := template.ParseGlob("templates/html/*.html")
	if err != nil {
		log.Fatalln(err)
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	tpl.ExecuteTemplate(res, "register.html", nil)
}
