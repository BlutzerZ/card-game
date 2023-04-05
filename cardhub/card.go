package cardhub

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func generate_room_id() string {
	rand.Seed(time.Now().UnixNano())

	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var result string
	for i := 0; i < 6; i++ {
		result += string(letters[rand.Intn(len(letters))])
	}

	return result
}

var connections = make(map[*websocket.Conn]string)

type Player struct {
	ID         string
	Connection *websocket.Conn
}

type Room struct {
	ID     string
	Player []Player
}

var rooms []Room

func GameWS(c *gin.Context) {
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer ws.Close()

	// generate id of player
	var player Player
	player.ID = uuid.New().String()
	player.Connection = ws

	// Check room id of url
	var room Room
	roomID := c.Param("roomID")
	if roomID == "" {
		// create new room
		room.ID = generate_room_id()
		room.Player = append(room.Player, player)

		rooms = append(rooms, room)
	} else {
		for _, room := range rooms {
			if roomID == room.ID {
				// add player to room
				room.Player = append(room.Player, player)
			}
		}
	}

	for {
		messageType, p, err := ws.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}

		for conn := range connections {
			err = conn.WriteMessage(messageType, p)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}
