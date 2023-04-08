package cardhub

import (
	"encoding/json"
	"fmt"

	"github.com/gorilla/websocket"
)

func StartGame(room Room) (Room, error) {
	// set playing mode to true
	room.isPlaying = true
	for i, playerReceiver := range room.Player {
		// get random card to player
		room.Player[i].Deck = GetRandomCard(7)
		// fmt.Println(playerReceiver.Deck)

		deckJSON, err := json.Marshal(playerReceiver.Deck)
		if err != nil {
			fmt.Println(err)
			return room, err
		}
		err = playerReceiver.Connection.WriteMessage(websocket.TextMessage, deckJSON)
		if err != nil {
			fmt.Println(err)
			return room, err
		}

	}
	fmt.Println(room)
	fmt.Println("=======")

	return room, nil
}

func DrawCard() {

}
