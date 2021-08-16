package usercontroller

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
	"webApp/auth"
	"webApp/database/model"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type userController struct {
	userCollection *mongo.Collection
}

func New(cl *mongo.Client) *userController {
	db := os.Getenv("DB_NAME")
	clt := os.Getenv("DB_USER_CLT")
	return &userController{
		userCollection: cl.Database(db).Collection(clt),
	}
}

func (uc *userController) GetUser(res http.ResponseWriter, req *http.Request) {
	log.Println("GetUser/:_id handler called")

	_id := mux.Vars(req)["_id"]
	id, err := primitive.ObjectIDFromHex(_id)
	if err != nil {
		log.Println("Error parsing objectID:", err)
		http.Error(res, "Error in request id", http.StatusBadRequest)
		return
	}
	log.Println("Searching for id:", id)

	var result model.User
	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	err = uc.userCollection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		log.Println("Error finding document:", err)
		http.Error(res, "Error finding user", http.StatusBadRequest)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	json.NewEncoder(res).Encode(result)
}

func (uc *userController) GetUsers(res http.ResponseWriter, _ *http.Request) {
	log.Println("GetUsers handler called")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cur, err := uc.userCollection.Find(ctx, bson.D{})
	if err != nil {
		log.Println("Error retrieving users:", err)
		http.Error(res, "Error retrieving users", http.StatusInternalServerError)
		return
	}

	var users []model.User

	for cur.Next(ctx) {
		var result model.User
		err := cur.Decode(&result)
		if err != nil {
			log.Println("Error decoding user:", err)
			http.Error(res, "Error iterating over users", http.StatusInternalServerError)
			return
		}
		fmt.Println(result)
		users = append(users, result)
	}
	res.Header().Set("Content-Type", "application/json")
	json.NewEncoder(res).Encode(users)
}

func (uc *userController) CreateUser(res http.ResponseWriter, req *http.Request) {
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

	err = auth.SignUp(u, uc.userCollection)
	if err != nil {
		log.Println("Error during sign-up:", err)
		http.Error(res, "Error signing up", http.StatusPermanentRedirect)
		return
	}

	// c := &http.Cookie{
	// 	Name:  os.Getenv("COOKIE_SID"),
	// 	Value: sID,
	// }
	// http.SetCookie(res, c)
	// Redirects user to login page
	http.Redirect(res, req, "/login", http.StatusTemporaryRedirect)
}

func (uc *userController) UpdateUser(res http.ResponseWriter, req *http.Request) {
	log.Println("UpdateUser/:id handler called")

	id := mux.Vars(req)["_id"]
	log.Println("Seraching for id:", id)

	var userInfo map[string]interface{}
	err := json.NewDecoder(req.Body).Decode(&userInfo)
	if err != nil {
		log.Println("Error decoding:", err)
		http.Error(res, "Error parsing body", http.StatusInternalServerError)
		return
	}
	fmt.Println(userInfo)

	filterId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println("Error parsing objectID:", err)
		http.Error(res, "Error in request id", http.StatusBadRequest)
		return
	}
	for key, val := range userInfo {
		u := bson.D{
			{"$set", bson.D{{key, val}}},
		}
		result, err := uc.userCollection.UpdateByID(context.Background(), filterId, u)
		if err != nil {
			log.Fatal("Error:", err)
			return
		}
		log.Println("MOdified:", result.ModifiedCount)
	}
}

func (uc *userController) DeleteUser(res http.ResponseWriter, req *http.Request) {
	log.Println("DeleteUser/:id handler called")

	id := mux.Vars(req)["_id"]
	log.Println("Seraching for id:", id)

	filterId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println("Error parsing objectID:", err)
		http.Error(res, "Error in request id", http.StatusBadRequest)
		return
	}

	result, err := uc.userCollection.DeleteOne(context.Background(), filterId)
	if err != nil {
		log.Fatal("Error:", err)
		return
	}
	log.Println("Modified:", result.DeletedCount)
}
