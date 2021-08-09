package main

import (
	"context"
	"log"
	"webApp/database/mongo_client"
	"webApp/router"
	"webApp/routes"

	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load()
}

func main() {
	// Connect to database
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

	r := router.NewRouter()
	if err = routes.AppendRoutes(r, m_cl); err != nil {
		log.Fatalln("Error appending routes:", err)
	}

}
