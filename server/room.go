package server

import (
	"encoding/json"
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

type Site struct {
	Index    int
	Name     string
	Sess     *Session
	Cards    card.Cards
	Rob      bool
	Status   int
}

func (s *Site) bind(sess *Session) {
	s.Name = sess.GetID()
	s.Sess = sess
}

func (s *Site) unbind(sess *Session) {
	s.Name = ""
	s.Sess = nil
}

type Room struct {
	Name string
	// 位置
	Sites map[int]*Site
	// 房价玩家
	Sessions map[string]*Session
	// 房间状态 0 未开始, 1 发牌中 , 2 进行中, 3 一局结束
	Status int

	// 底牌
	Cards card.Cards

	// 地主
	Landlord int
	//
	Owner string
	//
	Wait    int
	Created time.Time
	sync.Locker `json:"-"`
}

func NewRoom(name,owner string) *Room {
	r := &Room{
		Name:     name,
		Sessions: make(map[string]*Session,0),
		Sites:    make(map[int]*Site,3),

		Cards:   make(card.Cards,0),
		Owner:   owner,
		Created: time.Now(),
		Locker:  &sync.Mutex{},
	}

	r.Sites[0] = &Site{Index: 0}
	r.Sites[1] = &Site{Index: 1}
	r.Sites[2] = &Site{Index: 2}

	return r
}

// lock .
func (r *Room) lock() func(){
	r.Locker.Lock()
	return func() {
		r.Locker.Unlock()
	}
}

func (r *Room) start() error {
	// 洗牌
	cards := []string{"3","3","3","3","4","4","4","4","5","5","5","5","6","6","6","6","7","7","7","7","8","8","8","8","9","9","9","9","t","t","t","t","j","j","j","j","q","q","q","q","k","k","k","k","a","a","a","a","2","2","2","2","s","b"}
	for i := 53; i > 0; i-- {
		num := rand.Intn(i + 1)
		cards[i], cards[num] = cards[num], cards[i]
	}
	// 发牌
	index := 0
	for index < 51 {
		s := index%3
		r.Sites[s].Cards = append(r.Sites[s].Cards,card.NewCard(cards[index]))
		index++
	}

	// 保留底牌
	r.Cards = card.Cards{card.NewCard(cards[51]),card.NewCard(cards[52]),card.NewCard(cards[53])}
	r.Status = RoomStatusStart

	// 通知
	for _,s := range r.Sites {
		sort.Sort(s.Cards)
		s.Sess.conn.SendMessage([]byte("游戏开始,您的牌是:\n"+s.Cards.String()+"\n抢地主请在1分钟内输入[Rob]\n"))
	}

	// 等待抢地主
	go func() {
		time.Sleep(time.Second * 60)
		robs := make([]int,0)
		for _,s := range r.Sites {
			if s.Rob {
				robs = append(robs,s.Index)
			}
		}
		// 没有地主随机指定一位
		if len(robs) == 0 {
			r.Landlord = rand.Intn(3)
		} else {
			r.Landlord = robs[rand.Intn(len(robs))]
		}
		r.boardcast(fmt.Sprintf("地主是 %s",r.Sites[r.Landlord].Name))
		r.Sites[r.Landlord].Cards = append(r.Sites[r.Landlord].Cards,r.Cards...)
		sort.Sort(r.Sites[r.Landlord].Cards)
		r.Sites[r.Landlord].Sess.conn.SendMessage([]byte("您抢到了地主,您的牌是:\n"+r.Sites[r.Landlord].Cards.String()+"\n请出牌...\n"))

		// 游戏开始
		r.Status = RoomStatusPlaing
		r.Wait = r.Landlord
	}()
	return nil
}

func (r *Room) boardcast(msg string) {
	for _,s := range r.Sessions {
		s.conn.SendMessage([]byte(msg))
	}
}

// Join .
func (r *Room) Join(sess *Session) error{
	defer r.lock()()
	r.Sessions[sess.GetID()] = sess
	r.boardcast(fmt.Sprintf("玩家 %s 加入房间\n",sess.GetID()))
	for _,v := range r.Sites {
		if v.Sess == nil {
			v.bind(sess)
			break
		}
	}

	return nil
}

// OnDisconnect .
func (r *Room) Exit(sess *Session) {
	defer r.lock()()
	r.boardcast(fmt.Sprintf("玩家 %s 离开房间\n",sess.GetID()))
	delete(r.Sessions, sess.GetID())
	for _,v := range r.Sites {
		if v.Name == sess.GetID() {
			v.unbind(sess)
			break
		}
	}
}

// OnHandle .
func (r *Room) Handle(sess *Session,msg string) (err error) {
	defer r.lock()()
	if msg == "debug" {
		m,_ := json.Marshal(r)
		sess.conn.SendMessage(m)
		return
	}
	switch r.Status {
	case RoomStatusStart:
		if msg == "rob" {
			for _,v := range r.Sites {
				if v.Name == sess.GetID() {
					v.Rob = true
				}
			}
		}
	case RoomStatusPlaing:

	case RoomStatusInit:
		if r.Owner == sess.GetID() && msg == "start" {
			return r.start()
		}
	}

	err = errors.New(fmt.Sprintf("未知的指令:%s",msg))
	return 
}
