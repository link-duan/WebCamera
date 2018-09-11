package websocket

import (
	"github.com/Go-zh/net/websocket"
)

func response(ws *websocket.Conn, option string, message string) error {
	response := `{"option":"`+ option +`","message":"`+ message +`"}`
	return websocket.Message.Send(ws, response)
}