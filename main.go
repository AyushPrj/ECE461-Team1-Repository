package main

import (
	"ECE461-Team1-Repository/configs"
	// "bytes"
	"fmt"
	// "io/ioutil"
	"net/http"
	"os"
	// "path"
	// "strings"

	sw "ECE461-Team1-Repository/routes"
	templog "log"
)

func main() {
    // Run database
    configs.ConnectDB()
    templog.Printf("Server started")
    router := sw.NewRouter()

    // Serve the React app's static files
    // fs := CustomFileServer(http.Dir("./static/"))
    router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./assets/build/static/"))))
    router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, "./assets/build/index.html")
    })

    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    templog.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), CORSHandler(router)))
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
