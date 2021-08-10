package routes

import (
	"webApp/router"
	"webApp/routes/controllers"
	sessioncontroller "webApp/routes/controllers/sessionController"
	usercontroller "webApp/routes/controllers/userController"

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
	r.GET("/user/{_id}", uc.GetUser)
	r.GET("/users", uc.GetUsers)
	r.PUT("/user/{_id}", uc.UpdateUser)
	r.DELETE("/user/{_id}", uc.DeleteUser)

	return nil
}
