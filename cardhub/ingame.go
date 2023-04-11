package cardhub

import (
	"encoding/json"
	"fmt"
	"strings"

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

		// send current deck to player
		message := []byte(fmt.Sprintf("[Current-Card]: %s", room.Game.CurrentCard))
		err = playerReceiver.Connection.WriteMessage(websocket.TextMessage, message)
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
				// card is match the current card
				splitedCard := strings.Split(card[1], "_")

				// ============ for testing only ====================
				if splitedCard[1] == "reverse" {
					fmt.Println("reversed")
					for i, j := 0, len(room.Game.Queue.Player)-1; i < j; i, j = i+1, j-1 {
						room.Game.Queue.Player[i], room.Game.Queue.Player[j] = room.Game.Queue.Player[j], room.Game.Queue.Player[i]
					}
					// matching index
					for ip, p := range room.Game.Queue.Player {
						if player.ID == p.ID {
							if ip+1 >= len(room.Game.Queue.Player)-1 {
								room.Game.Queue.Turn = 0
							} else {
								room.Game.Queue.Turn = ip + 1
							}
							break
						}
					}

					for _, p := range room.Player {
						message := []byte("[Current-Queue]: Reversed")
						err := p.Connection.WriteMessage(websocket.TextMessage, message)
						if err != nil {
							fmt.Println(err)
							return room, err
						}
					}

				} else if splitedCard[1] == "skip" {
					if room.Game.Queue.Turn+2 == len(room.Game.Queue.Player) {
						room.Game.Queue.Turn = 0
					} else if room.Game.Queue.Turn+2 > len(room.Game.Queue.Player) {
						room.Game.Queue.Turn = 1
					} else {
						room.Game.Queue.Turn += 2
					}

				} else {
					if room.Game.Queue.Turn+1 >= len(room.Game.Queue.Player) {
						room.Game.Queue.Turn = 0
					} else {
						room.Game.Queue.Turn += 1
					}
				}
				break
				//  ===== for testing only =======

				if card[1] == deckCard && (strings.Contains(string(room.Game.CurrentCard), splitedCard[0]) || strings.Contains(string(room.Game.CurrentCard), splitedCard[1])) {
					room.Player[ip].Deck = append(room.Player[ip].Deck[:i], room.Player[ip].Deck[i+1:]...) // remove card
					deckJSON, err := json.Marshal(room.Player[ip].Deck)
					if err != nil {
						return room, err
					}

					// update current card
					room.Game.CurrentCard = card[1]
					// queue next
					if splitedCard[1] == "reverse" {
						for i, j := 0, len(room.Game.Queue.Player)-1; i < j; i, j = i+1, j-1 {
							room.Game.Queue.Player[i], room.Game.Queue.Player[j] = room.Game.Queue.Player[j], room.Game.Queue.Player[i]
						}
						// matching index
						for ip, p := range room.Game.Queue.Player {
							if player.ID == p.ID {
								if room.Game.Queue.Turn+1 >= len(room.Game.Queue.Player) {
									room.Game.Queue.Turn = 0
								} else {
									room.Game.Queue.Turn = ip + 1
								}
								break
							}
						}

						for _, p := range room.Player {
							message := []byte("[Current-Queue]: Reversed")
							err := p.Connection.WriteMessage(websocket.TextMessage, message)
							if err != nil {
								fmt.Println(err)
								return room, err
							}
						}

					} else if splitedCard[1] == "skip" {
						if room.Game.Queue.Turn+2 == len(room.Game.Queue.Player) {
							room.Game.Queue.Turn = 0
						} else if room.Game.Queue.Turn+2 > len(room.Game.Queue.Player) {
							room.Game.Queue.Turn = 1
						} else {
							room.Game.Queue.Turn += 2
						}

					} else {
						if room.Game.Queue.Turn+1 >= len(room.Game.Queue.Player) {
							room.Game.Queue.Turn = 0
						} else {
							room.Game.Queue.Turn += 1
						}
					}

					player.Connection.WriteMessage(websocket.TextMessage, deckJSON)

					// Notify all player the currnent card
					for _, p := range room.Player {
						message := []byte(fmt.Sprintf("[Current-Card]: %s", room.Game.CurrentCard))
						err := p.Connection.WriteMessage(websocket.TextMessage, message)
						if err != nil {
							fmt.Println(err)
							return room, err
						}
					}
					break
				}
			}
		}
	}
	return room, nil
}

func TakeCard(player Player, room Room) (Room, error) {
	for ip, p := range room.Player {
		if player.ID == p.ID && player.ID == room.Game.Queue.Player[room.Game.Queue.Turn].ID {
			room.Player[ip].Deck = append(room.Player[ip].Deck, GetRandomCard(1)...)
			// send all card with new card
			deckJSON, err := json.Marshal(room.Player[ip].Deck)
			if err != nil {
				fmt.Println(err)
				return room, err
			}
			err = p.Connection.WriteMessage(websocket.TextMessage, deckJSON)
			if err != nil {
				fmt.Println(err)
				return room, err
			}
			// queue next
			if room.Game.Queue.Turn+1 == len(room.Game.Queue.Player) {
				room.Game.Queue.Turn = 0
			} else {
				room.Game.Queue.Turn += 1
			}

		}
	}
	return room, nil
}
