package chatserver

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/matter123/chat/token"
	"time"
)

//UserConn is a websocket connection of a chat user
type UserConn struct {
	Username   string
	Token      string
	Connection websocket.Conn
	Send       chan interface{}
}

//Close closes the websocket from the server and and renews the Token a final time
func (u *UserConn) Close() {
	token.Renew(u.Token)
	close(u.Send)
	u.Connection.Close()
}

//todo move to own file
type PingMessage struct {
	Msgtype string
	Time    string
}

func (u *UserConn) Pump() {
	ping := time.NewTicker(45)
	defer func() {
		ping.Stop()
		u.Close()
	}()
	recievedPong := true
	for {
		select {
		case <-ping.C:
			if recievedPong == false {
				u.Close()
				return
			}
			u.Send <- PingMessage{
				Msgtype: "ping",
				Time:    string(time.Now().Unix()),
			}
			recievedPong = false
		case v := <-u.Send:
			err := u.Connection.WriteMessage(websocket.TextMessage, json.Marshal(v))
			if err != nil {
				u.Close()
				return
			}
		default:

		}
	}
}

//Broadcast is a pool of UserConn why would like to send and recieve broadcast
type Broadcast struct {
	//is a map instead of slice because y syntax is convienient for leave
	conns     map[*UserConn]bool
	broadcast chan interface{}
	join      chan *UserConn
	leave     chan *UserConn
}

//Join adds the UserConn to the broadcast pool, blocking until done so
func (b *Broadcast) Join(conn *UserConn) {
	b.join <- conn
}

//Leave removes the UserConn from the broadcast pool, blocking until done so
func (b *Broadcast) Leave(conn *UserConn) {
	b.leave <- conn
}

//Broadcast sends a struct to each UserConn in the broadcast pool
func (b *Broadcast) Broadcast(v interface{}) {
	b.broadcast <- v
}

//Pump continuously polls the broacast pool and distributes messages
func (b *Broadcast) Pump() {
	for {
		select {
		case conn := <-b.join:
			b.conns[conn] = true
		case conn := <-b.leave:
			delete(b.conns, conn)
			conn.Close()
		case v := <-b.broadcast:
			for conn := range b.conns {
				select {
				case conn.Send <- v:
				default:
					conn.Close()
					delete(b.conns, conn)
				}
			}
		}
	}
}

var b = &Broadcast{
	conns:     make(map[*UserConn]bool),
	broadcast: make(chan interface{}, 25),
	join:      make(chan *UserConn, 5),
	leave:     make(chan *UserConn, 5),
}

//Pool returns a ponter to the global Broadcast Pool
func Pool() *Broadcast {
	return b
}
