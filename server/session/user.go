package session

import (
	"fmt"
	"github.com/df-mc/dragonfly/server/player"
	"strconv"
	"time"
)

type User struct {
	Player       *player.Player
	Disconnected chan bool
	Ticker       time.Ticker
}

var users = make(map[string]*User)

func CreateUser(p *player.Player) User {
	if _, ok := users[p.Name()]; ok {
		return *users[p.Name()]
	}
	u := User{Player: p, Disconnected: make(chan bool)}
	users[p.Name()] = &u
	return u
}

func GetUser(p *player.Player) (*User, error) {
	if _, ok := users[p.Name()]; !ok {
		return nil, fmt.Errorf("user %s not found", p.Name())
	}
	return users[p.Name()], nil
}

func RemoveUser(p *player.Player) {
	u := users[p.Name()]
	u.Disconnected <- true
	u.Ticker.Stop()
	delete(users, p.Name())
}

func (u *User) Initialize() {
	ticker := time.NewTicker(1 * time.Second)
	u.Ticker = *ticker
	go func() {
		for {
			select {
			case <-u.Disconnected:
				return
			case _ = <-ticker.C:
				p := u.Player
				x := strconv.Itoa(int(p.Position().X()))
				y := strconv.Itoa(int(p.Position().Y()))
				z := strconv.Itoa(int(p.Position().Z()))
				p.SendTip("Current pos: " + x + ":" + y + ":" + z + ":" + p.World().Name())
			}
		}
	}()
}
