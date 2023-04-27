package main

import (
	"ECE461-Team1-Repository/configs"
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"

	sw "ECE461-Team1-Repository/routes"
	templog "log"
)
type customFileServer struct {
    root http.FileSystem
}

func (c *customFileServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    upath := r.URL.Path
    if !strings.HasPrefix(upath, "/") {
        upath = "/" + upath
        r.URL.Path = upath
    }

    f, err := c.root.Open(upath)
    if err != nil {
        http.ServeFile(w, r, path.Join("static", "index.html"))
        return
    }
    defer f.Close()

    d, err := f.Stat()
    if err != nil {
        http.ServeFile(w, r, path.Join("static", "index.html"))
        return
    }

    // Set the correct MIME type for JavaScript and CSS files
    if strings.HasSuffix(upath, ".js") {
        w.Header().Set("Content-Type", "application/javascript")
    } else if strings.HasSuffix(upath, ".css") {
        w.Header().Set("Content-Type", "text/css")
    }

    content, err := ioutil.ReadAll(f)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    reader := bytes.NewReader(content)
    http.ServeContent(w, r, d.Name(), d.ModTime(), reader)
}

func CustomFileServer(root http.FileSystem) http.Handler {
    return &customFileServer{root}
}
func main() {
    // Run database
    configs.ConnectDB()
    templog.Printf("Server started")
    router := sw.NewRouter()

    // Serve the React app's static files
    router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

    // Catch-all route to serve the React app's index.html file for client-side routing
    router.PathPrefix("/").HandlerFunc(serveReactApp)

    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    templog.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), CORSHandler(router)))
}

func serveReactApp(w http.ResponseWriter, r *http.Request) {
    upath := r.URL.Path
    filePath := "static" + upath

    // Check if the file exists
    _, err := os.Stat(filePath)

    // If the file exists, set the correct MIME type
    if err == nil {
        if strings.HasSuffix(filePath, ".js") {
            w.Header().Set("Content-Type", "application/javascript")
        } else if strings.HasSuffix(filePath, ".css") {
            w.Header().Set("Content-Type", "text/css")
        }
        http.ServeFile(w, r, filePath)
    } else {
        http.ServeFile(w, r, "static/index.html")
    }
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
