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

	r := gin.Default()

	r.GET("/ws", func(c *gin.Context) {
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer conn.Close()

		id := uuid.New().String()
		err = conn.WriteMessage(websocket.TextMessage, []byte(id))
		if err != nil {
			fmt.Println(err)
			return
		}

		for {
			messageType, p, err := conn.ReadMessage()
			if err != nil {
				fmt.Println(err)
				return
			}

			fmt.Println(string(p))
			err = conn.WriteMessage(messageType, p)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	})

	r.Run(":8080")
}
