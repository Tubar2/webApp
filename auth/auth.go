package auth

import (
	"context"
	"errors"
	"log"
	"time"
	redisauth "webApp/auth/redisAuth"
	"webApp/database/model"
	"webApp/security"
	"webApp/utils/channels"

	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func SignUp(u *model.User, col *mongo.Collection, redisCl *redis.Client) (string, error) {
	log.Println("Signing up user:", u)
	done := make(chan bool)
	var err error

	go func(ch chan<- bool) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer close(ch)
		defer cancel()

		// Checks if other user has already chosen this username
		err = unameTaken(u.Username, col)
		if err != mongo.ErrNoDocuments {
			log.Println("Error creating new user. Uname already taken")
			err = errors.New("username already taken")
			ch <- false
			return
		}

		// Try to insert user in db
		_, err = col.InsertOne(ctx, u)
		if err != nil {
			log.Println("Error creating new user in mongoDb:", err)
			ch <- false
			return
		}

		ch <- true
	}(done)

	log.Println("Waiting response")
	if channels.OK(done) {
		log.Println("Good response")
		return redisauth.StartUserSession(u, redisCl)
	}
	log.Println("Bad response")

	return "", err
}

func Logout(redisCl *redis.Client, key string) (int64, error) {
	log.Println("Finding key:", key)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cmd := redisCl.Del(ctx, key)

	return cmd.Result()
}

func SignIn(uname, pw string, col *mongo.Collection, redisCl *redis.Client) (string, error) {
	log.Printf("Signing in {'uname': %s, 'pw': %s}", uname, pw)
	done := make(chan bool)
	var u model.User
	var err error

	go func(ch chan<- bool) {
		defer close(ch)

		filter := bson.D{primitive.E{Key: "username", Value: uname}}
		loginRes := col.FindOne(context.Background(), filter)
		err = loginRes.Err()

		if err == mongo.ErrNoDocuments {
			log.Println("User not registered:", err)
			ch <- false
			return
		}
		if err != nil {
			log.Println("Error finding user:", err)
			ch <- false
			return
		}

		loginRes.Decode(&u)

		err = security.Verify(u.Password, pw)
		if err != nil {
			ch <- false
			return
		}

		ch <- true
	}(done)

	if channels.OK(done) {
		return redisauth.StartUserSession(&u, redisCl)
	}

	return "", err
}

func unameTaken(uname string, col *mongo.Collection) error {
	filter := bson.D{primitive.E{Key: "username", Value: uname}}
	res := col.FindOne(context.Background(), filter)
	return res.Err()
}
