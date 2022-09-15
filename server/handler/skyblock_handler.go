package handler

import (
	"github.com/alvin0319/DFSkyBlock/server/session"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/event"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
)

type SkyBlockHandler struct {
	player.NopHandler
	p            *player.Player
	Disconnected chan bool
}

func NewHandler(p *player.Player) *SkyBlockHandler {
	return &SkyBlockHandler{p: p, Disconnected: make(chan bool)}
}

func (h *SkyBlockHandler) HandleChat(ctx *event.Context, text *string) {
	// TODO
}

func (h *SkyBlockHandler) HandleBlockBreak(ctx *event.Context, pos cube.Pos, drops *[]item.Stack) {
	// TODO
}

func (h *SkyBlockHandler) HandleQuit() {
	session.RemoveUser(h.p)
}
