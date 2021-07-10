package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/chat-golang/sockets"
	"github.com/gorilla/mux"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	router := mux.NewRouter()

	hub := sockets.CreateHub()
	go hub.RunHub()

	router.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		sockets.HandleConnections(w, r, hub)
	})

	path, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	router.PathPrefix("/").Handler(http.FileServer(http.Dir(path+"/public/")))

	fmt.Println("Server running on localhost:"+port)
	log.Fatal(http.ListenAndServe(":"+port, router))

	
}