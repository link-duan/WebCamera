package app

import (
	"net/http"
	SocketServices "github.com/yaphper/WebCamera/app/websocket"
	"golang.org/x/net/websocket"
)

func init() {
	http.Handle("/socket/video_stream", websocket.Handler(SocketServices.VideoStreamWS))

	http.Handle("/", http.FileServer(http.Dir("app/view")))
}

func StartServer() error {

	return http.ListenAndServe("127.0.0.1:9808", nil)
}