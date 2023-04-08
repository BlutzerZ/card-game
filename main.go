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

	r.Run(":8080")
}
