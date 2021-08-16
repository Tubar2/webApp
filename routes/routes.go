package routes

import (
	"net/http"
	"webApp/router"
	"webApp/routes/controllers"
	homecontroller "webApp/routes/controllers/homeController"
	sessioncontroller "webApp/routes/controllers/sessionController"
	tweetcontroller "webApp/routes/controllers/tweetController"
	usercontroller "webApp/routes/controllers/userController"
	authmiddleware "webApp/routes/middlewares/authMiddleware"

	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
)

func AppendRoutes(r *router.Router, m_cl *mongo.Client, r_cl *redis.Client) error {
	// RESOURCES
	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("./public"))))

	// MIDDLEWARES
	authM := authmiddleware.New(r_cl)

	// INDEX
	r.GET("/", controllers.Index)

	// HOME
	hc := homecontroller.New(r_cl)
	r.GET("/home", authM.AuthUser(hc.Home))

	// SESSION_API
	rc := sessioncontroller.New(m_cl, r_cl)
	r.GET("/login", rc.LoginPage)
	r.GET("/register", rc.RegisterPage)

	r.POST("/logout", rc.Logout)
	r.POST("/login", rc.Login)

	// r.POST("/register", rc.Register)

	// USER_API
	uc := usercontroller.New(m_cl) // Generates
	r.GET("/user/{_id}", authM.AuthUser(uc.GetUser))
	r.GET("/users", uc.GetUsers)
	r.POST("/user", uc.CreateUser)
	r.PUT("/user/{_id}", authM.AuthUser(uc.UpdateUser))
	r.DELETE("/user/{_id}", authM.AuthUser(uc.DeleteUser))

	tc := tweetcontroller.New(m_cl, r_cl)
	r.POST("/tweet", authM.AuthUser(tc.NewTweet))
	r.GET("/tweets", authM.AuthUser(tc.GetTweets))

	return nil
}
