package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

func main() {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	connections := make(map[*websocket.Conn]string)

	r := gin.Default()

	r.GET("/ws", func(c *gin.Context) {
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
				message := []byte("[" + id + "]: " + string(p))
				err = conn.WriteMessage(messageType, message)
				if err != nil {
					fmt.Println(err)
					return
				}
			}
		}
	})

	r.Run(":8080")
}
