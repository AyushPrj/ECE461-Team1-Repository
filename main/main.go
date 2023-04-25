package main

import (
	"fmt"
	templog "log"
	"net/http"
	"os"

	"ECE461-Team1-Repository/configs"
	sw "ECE461-Team1-Repository/routes"
)

func main() {
	// run database
	configs.ConnectDB()
	templog.Printf("Server started")
	router := sw.NewRouter()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	templog.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), CORSHandler(router)))
	// templog.Fatal(http.ListenAndServe("0.0.0.0:8080", CORSHandler(router)))
}

func CORSHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, PUT, PATCH, DELETE, GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, access-control-allow-origin, access-control-allow-headers, X-Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		h.ServeHTTP(w, r)
	})
}
