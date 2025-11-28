package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	db := ConnectDB()
	defer db.Close()

	SetupRoutes(db)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Server berjalan di port:", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
