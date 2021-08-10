package authmiddleware

import (
	"log"
	"net/http"
	"os"
	redisauth "webApp/auth/redisAuth"

	"github.com/go-redis/redis/v8"
)

type authMiddleware struct {
	redisClient *redis.Client
}

func New(rc *redis.Client) *authMiddleware {
	return &authMiddleware{
		redisClient: rc,
	}
}

func (am *authMiddleware) AuthUser(next http.HandlerFunc) func(http.ResponseWriter, *http.Request) {
	log.Println("Auth middleware called")

	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		cookie, err := req.Cookie(os.Getenv("COOKIE_SID"))
		if err == http.ErrNoCookie {
			log.Println("No Cookie")
			http.Redirect(res, req, "/login", http.StatusBadRequest)
			return
		}

		found, err := redisauth.ValidSession(cookie.Value, am.redisClient)
		if err != nil {
			http.Error(res, "Error retrieving session", http.StatusInternalServerError)
		}

		if found > 0 {
			log.Println("Valid session")
			next.ServeHTTP(res, req)
			return

		} else {
			res.Write([]byte("Invalid session"))
			log.Println("Invalid user session")
			return
		}

	})
}
