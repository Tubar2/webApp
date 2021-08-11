package routes

import (
	"webApp/router"
	"webApp/routes/controllers"
	sessioncontroller "webApp/routes/controllers/sessionController"
	tweetcontroller "webApp/routes/controllers/tweetController"
	usercontroller "webApp/routes/controllers/userController"
	authmiddleware "webApp/routes/middlewares/authMiddleware"

	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
)

func AppendRoutes(r *router.Router, m_cl *mongo.Client, r_cl *redis.Client) error {
	r.GET("/", controllers.Index)

	// SESSION_API
	rc := sessioncontroller.New(m_cl, r_cl)
	r.POST("/login", rc.Login)
	r.POST("/register", rc.Register)
	r.GET("/logout", rc.Logout)

	// USER_API
	uc := usercontroller.New(m_cl)
	authM := authmiddleware.New(r_cl)
	r.GET("/user/{_id}", authM.AuthUser(uc.GetUser))
	r.GET("/users", uc.GetUsers)
	r.PUT("/user/{_id}", authM.AuthUser(uc.UpdateUser))
	r.DELETE("/user/{_id}", authM.AuthUser(uc.DeleteUser))

	tc := tweetcontroller.New(m_cl, r_cl)
	r.POST("/tweet", authM.AuthUser(tc.NewTweet))
	r.GET("/tweets", authM.AuthUser(tc.GetTweets))

	return nil
}
