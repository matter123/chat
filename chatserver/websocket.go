package chatserver

import (
	j "encoding/json"
	"github.com/gorilla/websocket"
	"github.com/matter123/chat/chatserver/messages"
	"log"
	"net/http"
	"strconv"
	"time"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

//WebsocketHandlerMock is handle to /websocketmock, WebsocketHandlerMock fakes a chat server, clients should use this
//if they wish to test their client, WebsocketHandlerMock does not check cookies and assumes the user is test
func WebsocketHandlerMock(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("failed handshake: ", err.Error())
		return
	}
	conn.WriteJSON(messages.NewUserJoin("test"))
	ping := messages.NewPing()
	conn.WriteJSON(ping)
	json := &messages.Text{}
	for ; ; conn.ReadJSON(json) {
		l, _ := j.Marshal(json)
		log.Print(string(l))
		if json.Pong != nil {
			pingtime, _ := strconv.ParseInt(ping.Ping.NanoTime, 10, 64)
			pongtime, _ := strconv.ParseInt(json.Pong.NanoTime, 10, 64)
			if pingtime != pongtime {
				log.Print("failed ping nano check")
			}
			log.Print("ping nano time delta=", time.Now().UnixNano()-pingtime)
			break
		}
	}
	conn.WriteJSON(messages.NewUserJoin("test2"))
	conn.WriteJSON(messages.NewMessage("test2", "test2"))
	log.Print("message test")
	json = &messages.Text{}
	for ; ; conn.ReadJSON(json) {
		l, _ := j.Marshal(json)
		log.Print(string(l))
		if json.Pong != nil {
			pingtime, _ := strconv.ParseInt(ping.Ping.NanoTime, 10, 64)
			pongtime, _ := strconv.ParseInt(json.Pong.NanoTime, 10, 64)
			if pingtime != pongtime {
				log.Print("failed ping nano check")
			}
			log.Print("ping nano time delta=", time.Now().UnixNano()-pingtime)
			break
		}
		if json.Message != nil {
			if json.Message.User == "test" {
				if json.Message.Message != "test" {
					log.Print("failed message check")
				}
				conn.WriteJSON(ping)
			} else {
				log.Print("failed message check")
			}
		}
	}
	conn.WriteJSON(messages.NewUserLeave("test2"))
	conn.WriteJSON(messages.NewUserLeave("test"))
	conn.Close()
}
