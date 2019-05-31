package main

import (
	"net/http"
)

// getRoutes registers the routes and uses populateSession middleware for authentication
func getRoutes() {
	http.Handle("/v1/add/user", populateSession(http.HandlerFunc(addUserHandler)))
	http.Handle("/v1/remove/user", populateSession(http.HandlerFunc(removeUserHandler)))
	http.Handle("/v1/add/movie", populateSession(http.HandlerFunc(addMovieHandler)))
	http.Handle("/v1/remove/movie", populateSession(http.HandlerFunc(removeMovieHandler)))
	http.Handle("/v1/update/movie", populateSession(http.HandlerFunc(updateMovieHandler)))
	http.Handle("/v1/get/movie", populateSession(http.HandlerFunc(getMovieHandler)))
	return
}
