package server

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	// "github.com/gorilla/mux"
)

type Server struct {
	Router *mux.Router
}

func (server *Server) Initialize() {
	log.Println("Welcome to library")
	server.Router = mux.NewRouter()
	server.initializeRoute()
}

func (server *Server) Run(port string) {
	log.Printf("Server listening to Port %s\n", port)
	log.Fatal(http.ListenAndServe(port, server.Router))
}

func Run() {
	var server = Server{}

	server.Initialize()
	server.Run(":9000")
}
