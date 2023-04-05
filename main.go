package main

import (
	"card-game/cardhub"

	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.Default()

	cardhub.CreateCard()
	r.GET("/ws/:roomID", cardhub.GameWS)
	r.GET("/ws/", cardhub.GameWS)
	// ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// defer ws.Close()

	// id := uuid.New().String()

	// connections[ws] = id

	// for {
	// 	messageType, p, err := ws.ReadMessage()
	// 	if err != nil {
	// 		fmt.Println(err)
	// 		return
	// 	}

	// 	for conn := range connections {
	// 		message := []byte("[" + id + "]: " + string(p))
	// 		err = conn.WriteMessage(messageType, message)
	// 		if err != nil {
	// 			fmt.Println(err)
	// 			return
	// 		}
	// 	}
	// }

	r.Run(":8080")
}
