package main

import (
	"log"
	"net/http"
	"os"

	"go-xsd/internal/routes"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	routes.SetupRoutes(router)

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server running on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
