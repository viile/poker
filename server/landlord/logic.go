package landlord

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/viile/poker/tools/session"
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
	Index  int
	Name   string
	Sess   *session.Session
	Cards  Cards
	Rob    bool
	Status int
}

func (s *Site) bind(ctx context.Context, sess *session.Session) {
	s.Sess = sess
}

func (s *Site) unbind(ctx context.Context, sess *session.Session) {
	s.Name = ""
	s.Sess = nil
}

type Room struct {
	Id   int
	Name string
	// 位置
	Sites map[int]*Site
	// 玩家
	Sessions map[string]*session.Session
	// 房间状态 0 未开始, 1 发牌中 , 2 进行中, 3 一局结束
	Status int

	// 底牌
	Cards Cards

	// 地主
	Landlord int
	//
	Owner string
	// 上一手牌型
	LastCards Cards
	// 上一手牌位置
	LastIndex int
	// 等待出牌位置
	Wait int

	event chan string

	Created     time.Time
	sync.Locker `json:"-"`
}

func NewRoom(ctx context.Context, index int, name, owner string) *Room {
	r := &Room{
		Id:       index,
		Name:     name,
		Sessions: make(map[string]*session.Session, 0),
		Sites:    make(map[int]*Site, 3),
		Cards:    make(Cards, 0),
		event:    make(chan string, 1024),
		Owner:    owner,
		Created:  time.Now(),
		Locker:   &sync.Mutex{},
	}

	r.Sites[0] = &Site{Index: 0}
	r.Sites[1] = &Site{Index: 1}
	r.Sites[2] = &Site{Index: 2}

	return r
}

// lock .
func (r *Room) lock() func() {
	r.Locker.Lock()
	return func() {
		r.Locker.Unlock()
	}
}

func (r *Room) TypeName() string {
	return "斗地主"
}

func (r *Room) OnlineCounts() int {
	return len(r.Sessions)
}

func (r *Room) start(ctx context.Context) error {
	// 洗牌
	cards := []string{"3", "3", "3", "3", "4", "4", "4", "4", "5", "5", "5", "5", "6", "6", "6", "6", "7", "7", "7", "7", "8", "8", "8", "8", "9", "9", "9", "9", "t", "t", "t", "t", "j", "j", "j", "j", "q", "q", "q", "q", "k", "k", "k", "k", "a", "a", "a", "a", "2", "2", "2", "2", "s", "b"}
	for i := 53; i > 0; i-- {
		num := rand.Intn(i + 1)
		cards[i], cards[num] = cards[num], cards[i]
	}
	// 发牌
	r.Sites[0].Cards = nil
	r.Sites[1].Cards = nil
	r.Sites[2].Cards = nil
	index := 0
	for index < 51 {
		s := index % 3
		r.Sites[s].Cards = append(r.Sites[s].Cards, NewCard(cards[index]))
		index++
	}

	// 保留底牌
	r.Cards = Cards{NewCard(cards[51]), NewCard(cards[52]), NewCard(cards[53])}
	r.Status = RoomStatusStart

	// 通知
	for _, s := range r.Sites {
		sort.Sort(s.Cards)
		s.Sess.SendMsg("游戏开始,您的牌是:\n" + s.Cards.String() + "\n抢地主请在30秒内输入[rob]\n")
	}

	// 等待抢地主
	go func() {
		time.Sleep(time.Second * 30)
		robs := make([]int, 0)
		for _, s := range r.Sites {
			if s.Rob {
				robs = append(robs, s.Index)
			}
		}
		// 没有地主随机指定一位
		if len(robs) == 0 {
			r.Landlord = rand.Intn(3)
		} else {
			r.Landlord = robs[rand.Intn(len(robs))]
		}
		r.boardcast(fmt.Sprintf("地主是 %s\n", r.Sites[r.Landlord].Name))
		r.Sites[r.Landlord].Cards = append(r.Sites[r.Landlord].Cards, r.Cards...)
		sort.Sort(r.Sites[r.Landlord].Cards)
		r.Sites[r.Landlord].Sess.SendMsg("您抢到了地主,您的牌是:\n" + r.Sites[r.Landlord].Cards.String() + "\n请出牌...\n")

		// 游戏开始,重置状态
		r.Status = RoomStatusPlaing
		r.Wait = r.Landlord
		r.LastIndex = -1
		for _, v := range r.Sites {
			v.Rob = false
		}
	}()
	return nil
}

func (r *Room) changeWait() {
	r.Wait = (r.Wait + 1) % 3
	r.Sites[r.Wait].Sess.SendMsg(fmt.Sprintf("轮到您出牌了!\n"))
}

//
func (r *Room) pass(ctx context.Context, s *session.Session) (err error) {
	r.boardcast(fmt.Sprintf("玩家 %s 过牌\n", s.GetName(nil)))
	r.changeWait()
	if r.Wait == r.LastIndex {
		r.LastIndex = -1
		r.LastCards = nil
	}
	return
}

//
func (r *Room) calc(ctx context.Context, s *session.Session, e string) (err error) {

	c := NewCards(e)
	// 判断所出手牌是否合法牌型
	//var ct int
	if _, err = c.Parser(); err != nil {
		return
	}
	// 判断是否拥有改手牌
	if !r.Sites[r.Wait].Cards.Contain(c) {
		err = errors.New("出了不存在的牌,请重新选牌...")
		return
	}

	// 当前场面无上家出牌
	if r.LastIndex >= 0 && r.LastCards.Len() > 0 {
		var i int
		if i, err = r.LastCards.Compare(c); err != nil {
			return
		}
		if i != 1 {
			err = errors.New("出的牌比上家小,请重新选牌...")
			return
		}
	}

	r.Sites[r.Wait].Cards = r.Sites[r.Wait].Cards.Remove(c)

	r.boardcast(fmt.Sprintf("玩家 %s 出牌: %s\n", s.GetID(), c.String()))
	r.boardcast(fmt.Sprintf("玩家 %s 剩余手牌数: %d\n", s.GetID(), r.Sites[r.Wait].Cards.Len()))
	s.SendMsg(fmt.Sprintf("您的手牌:%s\n", r.Sites[r.Wait].Cards))
	if r.Sites[r.Wait].Cards.Len() <= 0 {
		r.Status = RoomStatusInit
		r.boardcast(fmt.Sprintf("游戏结束,玩家 %s 获胜\n", s.GetID()))
	}

	r.LastCards = c
	r.LastIndex = r.Wait

	r.changeWait()

	return
}

func (r *Room) boardcast(msg string) {
	for _, s := range r.Sessions {
		s.SendMsg(msg)
	}
}

// Join .
func (r *Room) Join(ctx context.Context, sess *session.Session) error {
	defer r.lock()()
	r.Sessions[sess.GetID()] = sess
	r.boardcast(fmt.Sprintf("玩家 %s 加入房间\n", sess.GetID()))
	for _, v := range r.Sites {
		if v.Sess == nil {
			v.bind(ctx, sess)
			break
		}
	}

	return nil
}

// Logout .
func (r *Room) Exit(ctx context.Context, sess *session.Session) {
	defer r.lock()()
	r.boardcast(fmt.Sprintf("玩家 %s 离开房间\n", sess.GetID()))
	delete(r.Sessions, sess.GetID())
	for _, v := range r.Sites {
		if v.Name == sess.GetID() {
			v.unbind(ctx, sess)
			break
		}
	}
}

// OnHandle .
func (r *Room) Handle(ctx context.Context, sess *session.Session, msg string) (err error) {
	defer r.lock()()
	if msg == "debug" {
		m, _ := json.Marshal(r)
		sess.SendMessage(m)
		return
	}
	switch r.Status {
	case RoomStatusStart:
		if msg == "rob" {
			for _, v := range r.Sites {
				if v.Name == sess.GetID() {
					v.Rob = true
					return
				}
			}
		}
	case RoomStatusPlaing:
		if r.Sites[r.Wait].Name != sess.GetID() {
			err = errors.New("还未轮到您出牌!")
			return
		}

		if msg == "pass" {
			return r.pass(ctx, sess)
		} else {
			return r.calc(ctx, sess, msg)
		}

	case RoomStatusInit:
		if r.Owner == sess.GetID() && msg == "start" {
			return r.start(ctx)
		}
	}

	err = errors.New(fmt.Sprintf("未知的指令:%s", msg))
	return
}
