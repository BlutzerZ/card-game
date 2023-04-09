package cardhub

import (
	"fmt"
	"math/rand"
	"net/http"
	"strings"
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

type Player struct {
	ID         string
	Deck       []string
	Connection *websocket.Conn
}

type PlayerQueue struct {
	Player []Player
	Turn   int
}

type OnGame struct {
	Queue       PlayerQueue
	CurrentCard string
}

type Room struct {
	ID        string
	Player    []Player
	isPlaying bool
	Game      OnGame
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
	fmt.Println(player.ID)

	// Check room id of url
	var room Room
	roomID := c.Param("roomID")
	if roomID == "" {
		// create new room
		room.isPlaying = false
		room.ID = generate_room_id()
		room.Player = append(room.Player, player)
		fmt.Println(room.ID)
		message := []byte(fmt.Sprintf("[NEW=ROOM] Your invite link is %s", room.ID))
		ws.WriteMessage(websocket.TextMessage, message)

		rooms = append(rooms, room)
	} else {
		isFound := false
		for i, r := range rooms {
			fmt.Println(roomID, room.ID)
			if roomID == r.ID {
				// cbeck if game is already start
				if r.isPlaying {
					ws.WriteMessage(websocket.TextMessage, []byte("game already start"))
					return
				}

				// add player to room
				r.Player = append(r.Player, player)
				isFound = true
				room = r
				rooms[i] = room
				break
			}
		}

		if !(isFound) {
			return
		}

	}

	for {
		// read message ==============================
		messageType, p, err := ws.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}

		// matching room =====================
		indexRoom := 0
		for i, r := range rooms {
			if room.ID == r.ID {
				room = r
				indexRoom = i
				break
			}
		}

		// check if leader start game
		if strings.Contains(string(p), "/start") && !room.isPlaying {
			room, err = StartGame(room)
			if err != nil {
				fmt.Println(err)
				return
			}
			// for _, p := range card
			for _, playerReceiver := range room.Player {
				fmt.Println(playerReceiver.Deck)
			}
			fmt.Println(room.Game.Queue.Player[room.Game.Queue.Turn].ID)
			rooms[indexRoom] = room
			continue
		}

		// if player throw card
		if strings.Contains(string(p), "/throw") {
			fmt.Println("try throwing card")

			card := strings.Split(string(p), " ")
			room, err = Throwcard(card, player, room)
			if err != nil {
				fmt.Println(err)
				return
			}
			rooms[indexRoom] = room
			continue
		}

		for _, playerReceiver := range room.Player {
			message := []byte(fmt.Sprintf("[%s]: %s", player.ID, string(p)))
			err = playerReceiver.Connection.WriteMessage(messageType, message)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}
