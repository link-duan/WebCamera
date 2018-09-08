package app

import (
	"net/http"
	"github.com/Go-zh/net/websocket"
	SocketServices "github.com/yaphper/WebCamera/app/websocket"
)

func init() {
	http.Handle("/socket/", websocket.Handler(SocketServices.VideoStreamWS))

	http.Handle("/", http.FileServer(http.Dir("app/view")))
}

func StartServer() error {

	return http.ListenAndServe(":9808", nil)
}