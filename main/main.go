package main

import (
	templog "log"

	"net/http"

	"github.com/gorilla/handlers"
	//"github.com/gorilla/mux"

	"ECE461-Team1-Repository/configs"

	sw "ECE461-Team1-Repository/routes"
)

func main() {
	//run database
	configs.ConnectDB()
	templog.Printf("Server started")
	router := sw.NewRouter()

	headersOk := handlers.AllowedHeaders([]string{"Content-Type", "access-control-allow-origin", "access-control-allow-headers"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"POST", "PUT", "PATCH", "DELETE"})

	
	templog.Fatal(http.ListenAndServe(":8080", handlers.CORS(originsOk, headersOk, methodsOk)(router)))
}
