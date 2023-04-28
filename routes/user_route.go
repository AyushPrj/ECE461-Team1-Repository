package routes

import (
	"fmt"
	"net/http"
	"strings"

	"ECE461-Team1-Repository/controllers"
	"github.com/gorilla/mux"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
		handler = Logger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	return router
}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!")
}

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		Index,
	},

	Route{
		"CreateAuthToken",
		strings.ToUpper("Put"),
		"/authenticate",
		controllers.CreateAuthToken,
	},

	Route{
		"PackageByNameDelete",
		strings.ToUpper("Delete"),
		"/package/byName/{name}",
		controllers.PackageByNameDelete,
	},

	Route{
		"PackageByNameGet",
		strings.ToUpper("Get"),
		"/package/byName/{name}",
		controllers.PackageByNameGet,
	},

	Route{
		"PackageByRegExGet",
		strings.ToUpper("Post"),
		"/package/byRegEx",
		controllers.PackageByRegExGet,
	},

	Route{
		"PackageCreate",
		strings.ToUpper("Post"),
		"/package",
		controllers.PackageCreate,
	},

	Route{
		"PackageDelete",
		strings.ToUpper("Delete"),
		"/package/{id}",
		controllers.PackageDelete,
	},

	Route{
		"PackageRate",
		strings.ToUpper("Get"),
		"/package/{id}/rate",
		controllers.PackageRate,
	},

	Route{
		"PackageRetrieve",
		strings.ToUpper("Get"),
		"/package/{id}",
		controllers.PackageRetrieve,
	},

	Route{
		"PackageUpdate",
		strings.ToUpper("Put"),
		"/package/{id}",
		controllers.PackageUpdate,
	},

	Route{
		"PackagesList",
		strings.ToUpper("Post"),
		"/packages",
		controllers.PackagesList,
	},

	Route{
		"RegistryReset",
		strings.ToUpper("Delete"),
		"/reset",
		controllers.RegistryReset,
	},
}
