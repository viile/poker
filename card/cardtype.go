package card

import "errors"

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

const (
	StateSole = iota
)

func ParserCards(str string) (t int,err error) {
	if str == "" {
		return 0, errors.New("长度错误")
	}
	l := len(str)
	if l == 1 {
		t = TypeSole
		return
	}
	if l == 2 && str == "sb" {
		t = TypeJokerBomb
		return
	}
	var status = TypeSole
	var index = 1
	var tmp = string(str[0])
	for index < l {
		switch status {
		case TypeSole:
			s := string(str[l])
			// pair
			if s == tmp {

			}
			// chain
			if CardName2Value(tmp) + 1 == CardName2Value(s) {

			}
		}
		index++
	}

	return status, nil
}