package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/mkapnick/go-mirror/server/auth"
	"github.com/mkapnick/go-mirror/server/html"
	"github.com/mkapnick/go-mirror/server/routes"
	"github.com/mkapnick/go-mirror/server/ws"
	"golang.org/x/net/websocket"
	"net/http"
)

func main() {
	// mux router
	r := mux.NewRouter()

	// root route, render the auth page
	r.HandleFunc("/", html.LandingHandler)

	// authenticate users and assign jwt token if valid
	r.HandleFunc("/auth/login", auth.AuthLogin).Methods("POST")

	// create user and assign jwt token if valid
	r.HandleFunc("/auth/create", auth.AuthCreate).Methods("POST")

	// get session user and info related to the session user
	r.HandleFunc("/me", auth.LoginFilter(routes.GetMe)).Methods("GET")

	// render chat template if logged in
	r.HandleFunc("/chat", auth.LoginFilter(routes.Chat)).Methods("GET")

	// handle websocket connection
	r.Handle("/ws", websocket.Handler(ws.HandleConnections))

	// public files
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./public/")))

	// register the mux
	http.Handle("/", r)

	fmt.Println("Starting server on port 9000")
	err := http.ListenAndServe(":9000", nil)

	if err != nil {
		fmt.Println("Error starting server on port 9000")
		panic(err)
	}
}
