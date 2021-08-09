package main

import (
	"context"
	"log"
	mongoCl "webApp/database/client"

	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load()
}

func main() {
	// Connect to database
	m_cl, err := mongoCl.NewClient()
	if err != nil {
		log.Fatalln("Error connecting to mongo database:", err)
	}
	defer func() {
		log.Println("Disconnecting from mongo database")
		if err = m_cl.Disconnect(context.Background()); err != nil {
			panic(err)
		}
	}()
}
