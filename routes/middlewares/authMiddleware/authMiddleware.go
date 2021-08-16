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

// AuthUser checks if cookie is present and if it holds a valid
// session
func (am *authMiddleware) AuthUser(next http.HandlerFunc) func(http.ResponseWriter, *http.Request) {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		log.Println("Auth middleware called")

		// Verify if cookie is present
		// if not, redirect to login
		cookie, err := req.Cookie(os.Getenv("COOKIE_SID"))
		if err == http.ErrNoCookie {
			log.Println("No Cookie")
			http.Redirect(res, req, "/login", http.StatusSeeOther)
			return
		}

		found, err := redisauth.ValidSession(cookie.Value, am.redisClient)
		if err != nil {
			http.Error(res, "Error retrieving session", http.StatusInternalServerError)
			http.Redirect(res, req, "/login", http.StatusSeeOther)
			return
		}

		if found > 0 {
			log.Println("Valid session")
			next.ServeHTTP(res, req)
			return

		} else {
			res.Write([]byte("Invalid session"))
			http.Redirect(res, req, "/login", http.StatusSeeOther)
			log.Println("Invalid user session")
			return
		}

	})
}
