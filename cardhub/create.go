package cardhub

import (
	"crypto/rand"
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

func GetRandomCard(totalCard int) []string {
	// var seed int64

	// // rand.Read(((*[8]byte)(unsafe.Pointer(&seed)))[:])
	// // rand.Seed(seed)

	var result []string
	for i := 0; i < totalCard; i++ {
		r := string(cardList[rand.Intn(len(cardList))])
		result = append(result, r)
	}

	return result
}
