package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
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

/* ############################################################
 Handlers
TODO: Add authentication wrapper
############################################################ */

func Index(w http.ResponseWriter, r *http.Request) {
	// FIXME: Authenticate
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Index")
}

// NewList gets a post and then inserts a json objects
// into that pgSQL db np
// TODO: add user who created it, via authentication
func NewList(w http.ResponseWriter, r *http.Request) {
	list := new(List)
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		fmt.Println(err) // TODO: Write Error to JSON
	}
	err = r.Body.Close()
	if err != nil {
		fmt.Println(err)
	}

	err = json.Unmarshal(body, list)
	if err != nil {
		w.Header().Set("Content-Type", "application/json;charset=UTF-8")
		w.WriteHeader(http.StatusUnprocessableEntity) //422
		err = json.NewEncoder(w).Encode(err)
		if err != nil {
			fmt.Println(err)
		}
	}

	// TODO: set list creator id
	_, err = InsertList(db, list) // NOTE: is the id automatically set?

	// Spit it back to the user
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.WriteHeader(http.StatusCreated) // 201?
	err = json.NewEncoder(w).Encode(list)
	if err != nil {
		fmt.Println(err)
	}
}

func NewCard(w http.ResponseWriter, r *http.Request) {
	card := new(Card)
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		fmt.Println(err)
	}
	err = r.Body.Close()
	if err != nil {
		fmt.Println(err)
	}

	err = json.Unmarshal(body, card)
	if err != nil {
		w.Header().Set("Content-Type", "application/json;charset=UTF-8")
		w.WriteHeader(http.StatusUnprocessableEntity) //422
		err = json.NewEncoder(w).Encode(err)
		if err != nil {
			fmt.Println(err)
		}
	}

	// TODO: set list creator id
	// TODO: Assign Members?
	_, err = InsertCard(db, card)

	// Spit it back to the user
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.WriteHeader(http.StatusCreated) // 201?
	err = json.NewEncoder(w).Encode(card)
	if err != nil {
		fmt.Println(err)
	}
}

func GetCard(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r) // GET a card by id
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		fmt.Println(err)
	}

	c, err := SelectCard(db, id)
	if err != nil {
		w.Header().Set("Content-Type", "application/json;charset=UTF-8")
		w.WriteHeader(http.StatusNotFound) // Doesn't exist
		err = json.NewEncoder(w).Encode(err)
		if err != nil {
			fmt.Println(err)
		}
	} else {
		w.Header().Set("Content-Type", "application/json;charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(c)
		if err != nil {
			fmt.Fprintf(w, "Error SON encoding %s", err)
		}
	}
}

func GetList(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r) // GET a card by id
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		fmt.Println(err)
	}

	l, err := SelectList(db, id)
	if err != nil {
		w.Header().Set("Content-Type", "application/json;charset=UTF-8")
		w.WriteHeader(http.StatusNotFound) // Doesn't exist
		err = json.NewEncoder(w).Encode(err)
		if err != nil {
			fmt.Println(err)
		}
	} else {
		w.Header().Set("Content-Type", "application/json;charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(l)
		if err != nil {
			fmt.Fprintf(w, "Error SON encoding %s", err)
		}
	}
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
