package landlord

import (
	"errors"
	"sort"
	"strings"
)

const (
	CardTwo        = "2"
	CardThree      = "3"
	CardFour       = "4"
	CardFive       = "5"
	CardSix        = "6"
	CardSeven      = "7"
	CardEight      = "8"
	CardNine       = "9"
	CardTen        = "t"
	CardJack       = "j"
	CardQueue      = "q"
	CardKing       = "k"
	CardAce        = "a"
	CardSmallJoker = "s"
	CardBigJoker   = "b"
)

func CardValue(n string) int {
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
	case CardBigJoker:
		return 17
	default:
		return 0
	}
}

type Card struct {
	Name string
	val  int
}

func NewCard(c string) Card {
	return Card{
		Name: c,
		val:  CardValue(c),
	}
}

func (c Card) GetName() string {
	return c.Name
}

func (c Card) GetVal() int {
	return c.val
}

const (
	// 单张
	TypeSole = iota
	// 顺子
	TypeSoleChain
	// 对子
	TypePair
	// 连队
	TypePairChain
	// 三张
	TypeTrio
	// 三带一
	TypeTrioSole
	// 三带二
	TypeTrioPair
	// 飞机
	TypeAirplane
	// 飞机带一
	TypeAirplaneSole
	// 飞机带二
	TypeAirplanePair
	// 四带二单张
	TypeDualSole
	// 四带二对子
	TypeDualPair
	// 炸弹
	TypeBomb
	// 王炸
	TypeJokerBomb
)

type Cards []Card

func NewCards(s string) Cards {
	cards := make(Cards, len(s))
	for k, v := range s {
		cards[k] = NewCard(string(v))
	}
	sort.Sort(cards)
	return cards
}

func (c Cards) Len() int           { return len(c) }
func (c Cards) Less(i, j int) bool { return c[i].GetVal() < c[j].GetVal() }
func (c Cards) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }

func (c Cards) String() string {
	var str string
	for _, s := range c {
		str += s.Name
	}
	return str
}
func (c Cards) IsSeq() bool {
	if c.Len() <= 0 {
		return false
	}

	head := c[0].GetVal()
	index := 0
	for index < c.Len() {
		if head+index != c[index].GetVal() {
			return false
		}
		index++
	}

	return true
}

func compare(a, b int) int {
	if a < b {
		return -1
	} else if a == b {
		return 0
	} else {
		return 1
	}
}

func (c Cards) Recommend(o Cards) (str string) {
	return ""
}

func (c Cards) Contain(o Cards) bool {
	index := 0
	flag := 0
	for index < c.Len() {
		if c[index].val == o[flag].val {
			flag++
		}
		if flag >= o.Len() {
			return true
		}
		index++
	}
	return false
}

// The result will be 0 if o==c, -1 if o < c, and +1 if o > c.
func (c Cards) Compare(o Cards) (i int, err error) {
	var a, b int
	if a, err = c.Parser(); err != nil {
		return
	}
	if b, err = o.Parser(); err != nil {
		return
	}
	// 先判断王炸
	if a == TypeJokerBomb {
		i = -1
		return
	}
	if b == TypeJokerBomb {
		i = 1
		return
	}

	// 再判断普通炸弹
	if a != b {
		if a == TypeBomb {
			i = -1
			return
		}
		if b == TypeBomb {
			i = 1
			return
		}

		// 如果双方牌型不一致,且都不是炸弹,说明牌型错误
		err = errors.New("不符合出牌规则")
		return
	}

	// 最后比较相同牌型大小
	if c.Len() != o.Len() {
		err = errors.New("不符合出牌规则")
		return
	}

	cc := NewCardsComponent(c)
	oc := NewCardsComponent(o)

	switch a {
	case TypeSole, TypeSoleChain, TypePair, TypePairChain, TypeTrio, TypeBomb:
		i = compare(o[0].GetVal(), c[0].GetVal())
		return
	case TypeTrioSole, TypeTrioPair, TypeAirplane, TypeAirplaneSole, TypeAirplanePair:
		i = compare(oc.Trio[0].GetVal(), cc.Trio[0].GetVal())
		return
	case TypeDualSole, TypeDualPair:
		i = compare(oc.Dual[0].GetVal(), cc.Dual[0].GetVal())
		return
	}

	err = errors.New("不符合出牌规则")
	return
}

func (c Cards) Remove(o Cards) Cards {
	ret := make(Cards, 0)
	var i, j int
	for i < c.Len() {
		if j >= o.Len() || c[i].Name != o[j].Name {
			ret = append(ret, c[i])
			i++
			continue
		}
		if c[i].Name == o[j].Name {
			i++
			j++
			continue
		}
	}
	sort.Sort(ret)
	return ret
}

func (c Cards) Parser() (t int, err error) {
	comp := NewCardsComponent(c)
	switch c.Len() {
	case 1:
		if comp.Sole.Len() == 1 {
			t = TypeSole
			return
		}
	case 2:
		// sb || 88
		if comp.Pair.Len() == 1 {
			t = TypePair
			return
		} else if c[0].GetName() == CardSmallJoker && c[1].GetName() == CardBigJoker {
			t = TypeJokerBomb
			return
		}
	case 3:
		// 888
		if comp.Trio.Len() == 1 {
			t = TypeTrio
			return
		}
	case 4:
		// 8888 || 888a
		if comp.Dual.Len() == 1 {
			t = TypeBomb
			return
		} else if comp.Sole.Len() == 1 && comp.Trio.Len() == 1 {
			t = TypeTrioSole
			return
		}
	case 5:
		// 888aa
		if comp.Pair.Len() == 1 && comp.Trio.Len() == 1 {
			t = TypeTrioPair
			return
		}
	case 6:
		// 348888
		if comp.Sole.Len() == 2 && comp.Dual.Len() == 1 {
			t = TypeDualSole
			return
		}
	case 8:
		// 33448888
		if comp.Pair.Len() == 2 && comp.Dual.Len() == 1 {
			t = TypeDualPair
			return
		}
	}

	// sole chain
	// 34567
	if c.Len() >= 5 && c[0].GetVal() <= 10 && c.Len() == comp.Sole.Len() && comp.Sole.IsSeq() {
		t = TypeSoleChain
		return
	}

	// pain chain
	// 33445566
	if c.Len() >= 6 && c.Len()%2 == 0 && c.Len() == comp.Pair.Len()*2 && comp.Pair.IsSeq() {
		t = TypePairChain
		return
	}

	// airTrio
	// 444555
	if c.Len() >= 6 && c.Len()%3 == 0 && c.Len() == comp.Trio.Len()*3 && comp.Trio.IsSeq() {
		t = TypeAirplane
		return
	}

	// airTrioSole
	// 34445556
	if c.Len() >= 8 && c.Len()%4 == 0 && (comp.Sole.Len()+comp.Pair.Len()*2+comp.Dual.Len()*4) == comp.Trio.Len() && comp.Trio.IsSeq() {
		t = TypeAirplaneSole
		return
	}

	// airTrioPair
	// 3344455566
	if c.Len() >= 10 && c.Len()%5 == 0 && (comp.Pair.Len()+comp.Dual.Len()*2) == comp.Trio.Len() && comp.Trio.IsSeq() {
		t = TypeAirplanePair
		return
	}

	err = errors.New("不符合出牌规则")
	return
}

type CardsComponent struct {
	Sole Cards
	Pair Cards
	Trio Cards
	Dual Cards
}

func NewCardsComponent(c Cards) *CardsComponent {
	comp := &CardsComponent{
		Sole: make(Cards, 0),
		Pair: make(Cards, 0),
		Trio: make(Cards, 0),
		Dual: make(Cards, 0),
	}
	tmp := make(map[string]int, 0)
	for _, v := range c {
		if _, ok := tmp[v.Name]; ok {
			tmp[v.Name]++
		} else {
			tmp[v.Name] = 1
		}
	}
	for k, v := range tmp {
		switch v {
		case 1:
			comp.Sole = append(comp.Sole, NewCard(k))
		case 2:
			comp.Pair = append(comp.Pair, NewCard(k))
		case 3:
			comp.Trio = append(comp.Trio, NewCard(k))
		case 4:
			comp.Dual = append(comp.Dual, NewCard(k))
		}
	}
	sort.Sort(comp.Sole)
	sort.Sort(comp.Pair)
	sort.Sort(comp.Trio)
	sort.Sort(comp.Dual)
	return comp
}
