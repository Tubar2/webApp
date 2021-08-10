package main

import (
	"context"
	"log"
	"net/http"
	redisauth "webApp/auth/redisAuth"
	"webApp/database/mongo_client"
	"webApp/router"
	"webApp/routes"

	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load()
}

func main() {
	// Connect to mongo database
	m_cl, err := mongo_client.NewClient()
	if err != nil {
		log.Fatalln("Error connecting to mongo database:", err)
	}
	defer func() {
		log.Println("Disconnecting from mongo database")
		if err = m_cl.Disconnect(context.Background()); err != nil {
			panic(err)
		}
	}()

	// Connect to redis
	r_cl, err := redisauth.NewClient()
	if err != nil {
		log.Fatal("Error connecting to redis:", err)
	}

	// Ceate a new server router and add routes
	r := router.NewRouter()
	if err = routes.AppendRoutes(r, m_cl, r_cl); err != nil {
		log.Fatalln("Error appending routes:", err)
	}

	log.Println("Starting server")

	// Start Server
	log.Fatal(http.ListenAndServe(":8080", r))
}
