package cardhub

import (
	"encoding/json"
	"fmt"

	"github.com/gorilla/websocket"
)

func StartGame(room Room) (Room, error) {
	// [SETUP ROOM ]
	room.isPlaying = true
	room.Game.Queue.Player = room.Player
	room.Game.Queue.Turn = 0
	room.Game.CurrentCard = GetRandomCard(1)[0]

	for i, playerReceiver := range room.Player {
		// get random card to player
		deck := GetRandomCard(7)
		room.Player[i].Deck = deck
		fmt.Println(room.Player[i].Deck)
		fmt.Println("=======")
		fmt.Println(i)

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

func Throwcard(card []string, player Player, room Room) (Room, error) {
	// fmt.Println(card)
	for ip, p := range room.Player {
		// find match id in room
		if player.ID == p.ID && player.ID == room.Game.Queue.Player[room.Game.Queue.Turn].ID {
			for i, deckCard := range p.Deck {
				// check card in deck or not
				if card[1] == deckCard {
					room.Player[ip].Deck = append(room.Player[ip].Deck[:i], room.Player[ip].Deck[i+1:]...) // remove card
					deckJSON, err := json.Marshal(room.Player[ip].Deck)
					if err != nil {
						return room, err
					}

					// queue next

					if room.Game.Queue.Turn+1 == len(room.Game.Queue.Player) {
						room.Game.Queue.Turn = 0
					} else {
						room.Game.Queue.Turn += 1
					}
					player.Connection.WriteMessage(websocket.TextMessage, deckJSON)
					break
				}
			}
		}
	}
	return room, nil
}
