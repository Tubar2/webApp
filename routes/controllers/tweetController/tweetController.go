package tweetcontroller

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
	"webApp/database/model"
	redisaux "webApp/utils/redisAux"

	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type tweetController struct {
	tweetCollection *mongo.Collection
	redisClient     *redis.Client
}

func New(cl *mongo.Client, rc *redis.Client) *tweetController {
	db := os.Getenv("DB_NAME")
	clt := os.Getenv("DB_TWEET_CLT")
	return &tweetController{
		tweetCollection: cl.Database(db).Collection(clt),
		redisClient:     rc,
	}
}

func (tc *tweetController) NewTweet(res http.ResponseWriter, req *http.Request) {
	log.Println("New tweet handler called")

	c_name := os.Getenv("COOKIE_SID")
	cookie, err := req.Cookie(c_name)
	if err == http.ErrNoCookie {
		log.Println("No cookie during new tweet call")
		http.Error(res, "No cookie", http.StatusUnprocessableEntity)
		return
	}

	req.ParseMultipartForm(0)

	tweet := req.FormValue("tweet")
	_id, err := redisaux.GetSessionID(cookie.Value, tc.redisClient)
	if err != nil {
		log.Println("Error getting user _id")
		http.Error(res, "No session", http.StatusBadRequest)
		return
	}

	tw := model.NewTweet(_id, tweet)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = tc.tweetCollection.InsertOne(ctx, tw)
	if err != nil {
		log.Println("Error creating new tweet in mongoDb:", err)
		log.Println("Error creating tweet")
		http.Error(res, "Internal eror during tweel publication", http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	json.NewEncoder(res).Encode(tw)
}

func (tc *tweetController) GetTweets(res http.ResponseWriter, req *http.Request) {
	log.Println("Get tweets handler called")

	c_name := os.Getenv("COOKIE_SID")
	cookie, err := req.Cookie(c_name)
	if err == http.ErrNoCookie {
		log.Println("No cookie during new tweet call")
		http.Error(res, "No cookie", http.StatusUnprocessableEntity)
		return
	}

	// Get _id from current session
	a_id, err := redisaux.GetSessionID(cookie.Value, tc.redisClient)
	if err != nil {
		log.Println("Error getting session _id")
		http.Error(res, "No session", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.D{primitive.E{Key: "a_id", Value: a_id}}
	cur, err := tc.tweetCollection.Find(ctx, filter)
	if err != nil {
		log.Println("Error retrieving tweets:", err)
		http.Error(res, "Error retrieving tweets", http.StatusInternalServerError)
		return
	}

	var tweets []model.Tweet

	for cur.Next(ctx) {
		var result model.Tweet
		err := cur.Decode(&result)
		if err != nil {
			log.Println("Error decoding tweet:", err)
			http.Error(res, "Error iterating over tweets", http.StatusInternalServerError)
			return
		}
		fmt.Println(result)
		tweets = append(tweets, result)
	}

	res.Header().Set("Content-Type", "application/json")
	json.NewEncoder(res).Encode(tweets)
}
