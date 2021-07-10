package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"time"

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

	go func ()  {
		for {
			fmt.Println(runtime.NumGoroutine() - 3)	
			time.Sleep(time.Second * 10)
		}		
	}()

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