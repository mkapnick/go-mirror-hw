package routes

import (
	"encoding/json"
	"fmt"
	"github.com/mkapnick/go-mirror-hw/server/auth"
	"net/http"
	"os"
)

func Chat(w http.ResponseWriter, r *http.Request) {
	// serve the static chat html file
	pwd, _ := os.Getwd()
	http.ServeFile(w, r, pwd+"/public/chat.html")
}

func GetMe(w http.ResponseWriter, r *http.Request) {
	sessionJson, err := json.Marshal(auth.SessionUser)
	if err != nil {
		fmt.Println("Error parsing session user")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(sessionJson)
}
