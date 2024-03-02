package main

import (
	"fmt"
	"go-MongoApi/router"
	"log"
	"net/http"
)

func main() {
	fmt.Println("Mongo server connected")
	fmt.Println("Server is getting started")
	r := router.Routes()
	log.Fatal(http.ListenAndServe(":8000", r))
	fmt.Println("Listening and serving at port 8000")
}
