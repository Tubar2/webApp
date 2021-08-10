package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
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

	// Ceate a new router and add routes
	r := router.NewRouter()
	if err = routes.AppendRoutes(r, m_cl, r_cl); err != nil {
		log.Fatalln("Error appending routes:", err)
	}

	// Create server
	srv := &http.Server{
		Addr:         "localhost:8080",
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r,
	}

	// Start Server on new routine
	go func() {
		log.Println("Starting server")
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	// Graceful Shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	<-c

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	log.Println("Shutting down")
	srv.Shutdown(ctx)
}
