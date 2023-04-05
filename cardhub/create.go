package cardhub

import (
	"fmt"
	"strconv"
)

var cardList = []string{"b", "b_2", "b_4"}

func CreateCard() {
	// red
	cardColor := []string{"red", "blue", "green", "yellow"}
	cardRnS := []string{"r", "s"}
	for _, c := range cardColor {
		for i := 0; i < 10; i++ {
			card := fmt.Sprintf(c + "_" + strconv.Itoa(i))
			cardList = append(cardList, card)
		}
		for _, rns := range cardRnS {
			card := fmt.Sprintf(c + "_" + rns)
			cardList = append(cardList, card)
		}
	}

	for _, card := range cardList {
		fmt.Println(card)
	}
}
