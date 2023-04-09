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
		deck := GetRandomCard(7)
		room.Player[i].Deck = deck
		fmt.Println(room.Player[i].Deck)
		fmt.Println("=======")

		fmt.Println(i)

		// fmt.Println(playerReceiver.Deck)

		deckJSON, err := json.Marshal(room.Player[i].Deck)
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

	for i, _ := range room.Player {
		fmt.Println(room.Player[i])
	}

	return room, nil
}

func DrawCard() {

}
