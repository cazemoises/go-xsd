package routes

import (
	"go-xsd/internal/app"

	"github.com/gorilla/mux"
)

func SetupRoutes(router *mux.Router) {
	router.HandleFunc("/validate", app.ValidateXMLHandler).Methods("POST")
	router.HandleFunc("/convert/xml-to-json", app.HandleXMLToJSON).Methods("POST")
	router.HandleFunc("/convert/json-to-xml", app.HandleJSONToXML).Methods("POST")
}
