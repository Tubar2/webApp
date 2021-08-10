package redisauth

import (
	"context"
	"log"
	"os"
	"time"
	"webApp/database/model"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

// Creates a connection to the redis db
func NewClient() (*redis.Client, error) {
	rcl := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("RD_ADDR"),
		Password: "",
		DB:       0,
	})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Pings redis db to check connection
	_, err := rcl.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	log.Println("Created new redis client connection")

	return rcl, nil
}

func StartUserSession(u *model.User, redisCl *redis.Client) (string, error) {
	log.Println("Starting session for user:", u)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	sID := uuid.NewString()

	uj := u.ToHSET()

	status := redisCl.HSet(ctx, sID, uj)
	if status.Err() != nil {
		log.Println("Error starting session. id:", sID, "status:", status)
		return "", status.Err()
	}
	log.Println()

	return sID, nil
}


