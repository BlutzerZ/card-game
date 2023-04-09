package cardhub

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
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

	// // r := rand.New(rand.NewSource(time.Now().UnixNano()))
	// for i := range cardList {
	// 	b := make([]byte, 8)
	// 	rand.Read(b)
	// 	seed := binary.BigEndian.Uint64(b)
	// 	r := rand.New(rand.NewSource(int64(seed) ^ time.Now().UnixNano()))
	// 	// Generate nilai acak untuk indeks pertukaran
	// 	j := r.Intn(len(cardList))
	// 	cardList[i], cardList[j] = cardList[j], cardList[i]
	// }

	newCards := make([]string, len(cardList))
	copy(newCards, cardList)

	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)
	r.Shuffle(len(newCards), func(i, j int) { newCards[i], newCards[j] = newCards[j], newCards[i] })

	resultCard := newCards[:7]

	return resultCard
}
