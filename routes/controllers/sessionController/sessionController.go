package sessioncontroller

import (
	"log"
	"net/http"
	"os"
	"webApp/auth"
	"webApp/database/model"

	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
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

func (rc *sessionController) Register(res http.ResponseWriter, req *http.Request) {
	log.Println("Register handler called")

	req.ParseMultipartForm(0)

	name := req.FormValue("name")
	uname := req.FormValue("username")
	email := req.FormValue("email")
	pw, err := bcrypt.GenerateFromPassword([]byte(req.FormValue("password")), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}

	u := model.NewUser(name, uname, email, string(pw))

	sID, err := auth.SignUp(u, rc.userCollection, rc.redisClient)
	if err != nil {
		log.Println("Error during sign-up:", err)
		http.Error(res, "Error signing up", http.StatusInternalServerError)
		return
	}

	c := &http.Cookie{
		Name:  os.Getenv("COOKIE_SID"),
		Value: sID,
	}
	http.SetCookie(res, c)

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
}
