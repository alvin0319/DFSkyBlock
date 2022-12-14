package command

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
)

// Great code https://github.com/Prim69/df-build-server/blob/main/game_mode.go

type GameMode struct {
	GameMode string `cmd:"gamemode"`
}

func (g GameMode) Run(s cmd.Source, o *cmd.Output) {
	var name string
	var mode world.GameMode
	switch g.GameMode {
	case "survival", "0", "s":
		name, mode = "survival", world.GameModeSurvival
	case "creative", "1", "c":
		name, mode = "creative", world.GameModeCreative
	case "adventure", "2", "a":
		name, mode = "adventure", world.GameModeAdventure
	case "spectator", "3", "sp":
		name, mode = "spectator", world.GameModeSpectator
	}

	s.(*player.Player).SetGameMode(mode)
	o.Printf("You have set your gamemode to %s", name)
}

// Allow ...
func (GameMode) Allow(s cmd.Source) bool {
	_, ok := s.(*player.Player)
	return ok
}

type gameMode string

// Type ...
func (gameMode) Type() string {
	return "GameMode"
}

// Options ...
func (gameMode) Options(cmd.Source) []string {
	return []string{
		"survival", "0", "s",
		"creative", "1", "c",
		"adventure", "2", "a",
		"spectator", "3", "sp",
	}
}
