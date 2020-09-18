package card

import "strings"

const (
	SuitHeart = iota
	SuitSpade
	SuitClub
	SuitDiamond
	SuitSmallJoker
	SuitBigJoker
)

const (
	CardTwo = "2"
	CardThree = "3"
	CardFour = "4"
	CardFive = "5"
	CardSix = "6"
	CardSeven = "7"
	CardEight = "8"
	CardNine = "9"
	CardTen = "t"
	CardJack = "j"
	CardQueue = "q"
	CardKing = "k"
	CardAce = "a"
	CardSmallJoker = "s"
	CardBigJoker = "b"
)

func SuitDisplay(s int) string {
	switch s {
	case SuitHeart:
		return "♥️"
	case SuitSpade:
		return "♠️"
	case SuitClub:
		return "♣️"
	case SuitDiamond:
		return "♦️"
	case SuitSmallJoker:
		return "🤡"
	default:
		return "🤠"
	}
}

func CardName2Value(n string) int {
	switch strings.ToLower(n) {
	case CardThree:
		return 3
	case CardFour:
		return 4
	case CardFive:
		return 5
	case CardSix:
		return 6
	case CardSeven:
		return 7
	case CardEight:
		return 8
	case CardNine:
		return 9
	case CardTen:
		return 10
	case CardJack:
		return 11
	case CardQueue:
		return 12
	case CardKing:
		return 13
	case CardAce:
		return 14
	case CardTwo:
		return 15
	case CardSmallJoker:
		return 16
	default:
		return 17
	}
}

type Card struct {
	name string
	suit int
}
type Cards []Card
func (c Cards) Len() int           { return len(c) }
func (c Cards) Less(i, j int) bool { return CardName2Value(c[i].name) < CardName2Value(c[j].name) }
func (c Cards) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }
