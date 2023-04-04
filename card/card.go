package 

import (
	"fmt"
	"net/http"main

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var connections = make(map[*websocket.Conn]string)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var cardList = []string{"b_2", "b_4"}

func CreateCard() {
	// red
	cardColor := []string{"red", "blue", "green", "yellow"}
	for i := 0; i < 10; i++ {
		for _, c := range cardColor {
			card := fmt.Sprintf(c + "_" + string(i))
			cardList = append(cardList, card)
		}
	}

	for _, card := range cardList {
		fmt.Printf(card)
	}
}

func GameWS(c *gin.Context) {
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer ws.Close()

	id := uuid.New().String()
	connections[ws] = id

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
