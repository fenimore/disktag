package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

func Serve() {
	// TODO: set connection of DB to global
	router := NewRouter()
	fmt.Println("Serving On Port: 7575")

	err := http.ListenAndServe(":7575", router)
	if err != nil { // FIXME: Log fatal?
		fmt.Println(err)
	}
}

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

		router.Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}
	return router
}

// Define handlers in handlers.go
var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		Index,
	},
}

/* Handlers */

/* ############################################################
 Index
############################################################ */

func Index(w http.ResponseWriter, r *http.Request) {
	// FIXME: Authenticate
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Index")
}

/* ############################################################
 Logger
############################################################ */
func Logger(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		inner.ServeHTTP(w, r)

		log.Printf(
			"%s\t%s\t%s\t%s",
			r.Method,
			r.RequestURI,
			name,
			time.Since(start),
		)
	})
}
