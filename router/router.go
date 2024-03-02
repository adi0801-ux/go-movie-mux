package router

import (
	"github.com/gorilla/mux"
	"go-MongoApi/controller"
)

func Routes() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/api/movies", controller.GetAllMovie).Methods("GET")
	router.HandleFunc("/api/movie", controller.CreateMovie).Methods("POST")
	router.HandleFunc("/api/movies/{id}", controller.MarkAsWatched).Methods("PUT")
	router.HandleFunc("/api/movies/{id}", controller.DeleteAMovie).Methods("DELETE")
	router.HandleFunc("/api/movies", controller.DeleteAllMovies).Methods("DELETE")
	return router
}
