package messages

import (
	"strconv"
	"time"
)

//Message is a substruct representing a Users message
type Message struct {
	User               string `json:"user"`
	Message            string `json:"message"`
	UnformattedMessage string `json:"raw_message"`
	Index              int    `json:"index"`
}

//PingPong is a substruct that is used for both ping and pong messages
type PingPong struct {
	NanoTime string `json:"nanotime"`
}

//JoinLeave is a substruct for when a user joins or leaves the chat
type JoinLeave struct {
	User string `json:"user"`
}

//Text is the struct that combines all the substructs together
type Text struct {
	Time    string     `json:"time"`
	Message *Message   `json:"message,omitempty"`
	Ping    *PingPong  `json:"ping,omitempty"`
	Pong    *PingPong  `json:"pong,omitempty"`
	Join    *JoinLeave `json:"join,omitempty"`
	Leave   *JoinLeave `json:"leave,omitempty"`
}

//NewPing generates a new Text of subtype Ping
func NewPing() *Text {
	return &Text{
		Time: strconv.FormatInt(time.Now().Unix(), 10),
		Ping: &PingPong{
			NanoTime: strconv.FormatInt(time.Now().UnixNano(), 10),
		},
	}
}

//NewUserJoin generates a new Text of subtype Join
func NewUserJoin(user string) *Text {
	return &Text{
		Time: strconv.FormatInt(time.Now().Unix(), 10),
		Join: &JoinLeave{
			User: user,
		},
	}
}

//NewUserLeave generates a new Text of subtype Leave
func NewUserLeave(user string) *Text {
	return &Text{
		Time: strconv.FormatInt(time.Now().Unix(), 10),
		Leave: &JoinLeave{
			User: user,
		},
	}
}

//NewMessage generates a new Text of subtype Message
func NewMessage(user string, message string) *Text {
	return &Text{
		Time: strconv.FormatInt(time.Now().Unix(), 10),
		Message: &Message{
			User:               user,
			Message:            message,
			UnformattedMessage: message,
		},
	}
}
