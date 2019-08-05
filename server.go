package main

import (
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/kpswamy540/verloop/controller"
)

//Application routes and server start here
func main() {
	apiController := &controller.APIController{}
	router := mux.NewRouter()
	router.HandleFunc("/test", apiController.Test).Methods("GET")
	router.HandleFunc("/get-question-id", apiController.GetQuestionID).Methods("POST","OPTIONS")

	router.HandleFunc("/repos", apiController.GetRepositories).Methods("POST", "OPTIONS")
	//Server port listening
	http.ListenAndServe(":8080", handlers.CORS(
		handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Content-Length"}),
		handlers.AllowedOrigins([]string{"*"}))(router))
}
