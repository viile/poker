package server

import (
	"errors"
	"fmt"
	"github.com/viile/poker/card"
	"math/rand"
	"sort"
	"sync"
	"time"
)

const (
	RoomStatusInit = iota
	RoomStatusStart
	RoomStatusPlaing
) 

type Room struct {
	name string
	sessions map[string]*Session
	// 房间状态 0 未开始, 1 发牌中 , 2 进行中
	status int
	// 托管玩家
	autos map[string]*Session
	// 玩家手牌
	cards map[string]card.Cards
	// 底牌
	procards card.Cards
	owner string
	wait string
	created time.Time
	sync.Locker
}

func NewRoom(name,owner string) *Room {
	return &Room{
		name :name,
		sessions:make(map[string]*Session,0),
		autos:make(map[string]*Session,0),
		cards:make(map[string]card.Cards,0),
		procards:make(card.Cards,0),
		owner: owner,
		created: time.Now(),
		Locker:&sync.Mutex{},
	}
}

// OnConnect .
func (r *Room) lock() func(){
	r.Locker.Lock()
	return func() {
		r.Locker.Unlock()
	}
}

func (r *Room) start() error {
	// 发牌
	cards := []string{"3","3","3","3","4","4","4","4","5","5","5","5","6","6","6","6","7","7","7","7","8","8","8","8","9","9","9","9","t","t","t","t","j","j","j","j","q","q","q","q","k","k","k","k","a","a","a","a","2","2","2","2","s","b"}
	for i := 53; i > 0; i-- {
		num := rand.Intn(i + 1)
		cards[i], cards[num] = cards[num], cards[i]
	}
	index := 0
	for k,_ := range r.sessions {
		tmp := make(card.Cards,0)
		for i:=0;i<17;i++{
			tmp = append(tmp,card.NewCard(cards[i + index*17]))
		}
		sort.Sort(tmp)
		r.cards[k] = tmp
		index++
	}
	// 底牌
	r.procards = card.Cards{card.NewCard(cards[51]),card.NewCard(cards[52]),card.NewCard(cards[53])}
	r.status = RoomStatusStart
	//
	for k,s := range r.sessions {
		s.conn.SendMessage([]byte("游戏开始,您的牌是:\n"+r.cards[k].String()+"\n抢地主输入[rob]\n"))
	}
	// 等待抢地主
	return nil
}

// OnConnect .
func (r *Room) Join(sess *Session) error{
	defer r.lock()()
	if len(r.sessions) >= 3 {
		return errors.New("房间已坐满!")
	}
	for _,s := range r.sessions {
		s.conn.SendMessage([]byte(fmt.Sprintf("玩家 %s 加入房间\n",sess.GetID())))
	}
	r.sessions[sess.GetID()] = sess

	if len(r.sessions) >= 3 {
		return r.start()
	}

	return nil
}

// OnDisconnect .
func (r *Room) Exit(sess *Session) {
	defer r.lock()()
	switch r.status {
	case RoomStatusInit:
		delete(r.sessions, sess.GetID())
		delete(r.autos, sess.GetID())
	case RoomStatusStart:
		delete(r.sessions, sess.GetID())
		r.autos[sess.GetID()] = sess
	case RoomStatusPlaing:
		delete(r.sessions, sess.GetID())
		r.autos[sess.GetID()] = sess
	}
}

// OnHandle .
func (r *Room) Handle(sess *Session,name string) {

}
