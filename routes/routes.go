package routes

import (
	"webApp/router"
	"webApp/routes/controllers"
	usercontroller "webApp/routes/controllers/userController"

	"go.mongodb.org/mongo-driver/mongo"
)

func AppendRoutes(r *router.Router, m_cl *mongo.Client) error {
	r.GET("/", controllers.Index)

	// USER_API
	uc := usercontroller.New(m_cl)
	r.GET("/user/{_id}", uc.GetUser)
	r.GET("/users", uc.GetUsers)
	r.PUT("/user/{_id}", uc.UpdateUser)
	r.DELETE("/user/{_id}", uc.DeleteUser)

	return nil
}
