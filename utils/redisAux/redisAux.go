package redisaux

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

func GetSessionID(sID string, rc *redis.Client) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Retrieves user's mongo ID
	res := rc.HGet(ctx, sID, "mID")

	return res.Result()
}
