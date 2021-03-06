package main

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
	"github.com/stretchr/objx"
)

type client struct {
	socket   *websocket.Conn
	send     chan *message
	room     *room
	userData objx.Map
}

func (c *client) read() {
	for {
		var msg *message
		if err := c.socket.ReadJSON(&msg); err == nil {
			msg.When = jsonTime{time.Now()}
			msg.Name = c.userData["name"].(string)
			if avatarURL, ok := c.userData["avatar_url"]; ok {
				msg.AvatarURL = avatarURL.(string)
			}
			c.room.forward <- msg
		} else {
			log.Println("エラーが発生しました：", err)
			break
		}
	}
	c.socket.Close()
}

func (c *client) write() {
	for msg := range c.send {
		if err := c.socket.WriteJSON(msg); err != nil {
			break
		}
	}
	c.socket.Close()
}
