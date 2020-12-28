package main

import (
	"log"
	"net/http"

	"github.com/sourjp/go-gophercises/cyoa/story"
)

func main() {
	log.Println("Server is running...")
	server := http.Server{
		Addr: "127.0.0.1:8080",
	}
	http.HandleFunc("/", story.ViewPageHandler)
	http.HandleFunc("/intro", story.ViewPageHandler)
	http.HandleFunc("/new-york", story.ViewPageHandler)
	http.HandleFunc("/debate", story.ViewPageHandler)
	http.HandleFunc("/sean-kelly", story.ViewPageHandler)
	http.HandleFunc("/mark-bates", story.ViewPageHandler)
	http.HandleFunc("/dever", story.ViewPageHandler)
	http.HandleFunc("/home", story.ViewPageHandler)
	server.ListenAndServe()
}
