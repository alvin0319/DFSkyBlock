package main

import (
	"fmt"
	"github.com/alvin0319/DFSkyBlock/server/handler"
	"github.com/alvin0319/DFSkyBlock/server/world"
	"github.com/df-mc/dragonfly/server"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/player/chat"
	"github.com/pelletier/go-toml"
	"github.com/sirupsen/logrus"
	"os"
	"strconv"
	"time"
)

func main() {
	config, err := readConfig()
	log := logrus.New()
	log.Formatter = &logrus.TextFormatter{ForceColors: true}
	log.Level = logrus.InfoLevel
	chat.Global.Subscribe(chat.StdoutSubscriber{})
	if err != nil {
		panic(err)
	}
	srv := server.New(&config.Config, log)
	srv.CloseOnProgramEnd()

	skyBlockPath := "skyblock"
	if _, err := os.Stat(skyBlockPath); os.IsNotExist(err) {
		if err := os.Mkdir(skyBlockPath, 0755); err != nil {
			panic(err)
		}
	}

	var treeType world.TreeInterface

	switch config.SkyBlock.TreeType {
	case "oak":
		treeType = world.NewOakTree()
	case "birch":
		treeType = world.NewBirchTree(false)
	case "spruce":
		treeType = world.NewSpruceTree()
	case "jungle":
		treeType = world.NewJungleTree()
	default:
		treeType = world.NewOakTree()
	}

	log.Infof("Selected tree type: %v", treeType.GetName())

	worldManager := world.NewManager(treeType)

	newWorld, err := worldManager.NewWorld("test")
	if err != nil {
		panic(err)
	}

	log.Printf("Loaded custom world %v", newWorld.Name())

	if err := srv.Start(); err != nil {
		panic(err)
	}
	for srv.Accept(func(p *player.Player) {
		go func() {
			for {
				x := strconv.Itoa(int(p.Position().X()))
				y := strconv.Itoa(int(p.Position().Y()))
				z := strconv.Itoa(int(p.Position().Z()))
				p.SendTip("Current pos: " + x + ":" + y + ":" + z)
				time.Sleep(1 * time.Second)
			}
		}()
		newWorld.AddEntity(p)
		p.Handle(handler.NewHandler(p))
	}) {
	}
}

// readConfig reads the configuration from the config.toml file, or creates the file if it does not yet exist.
func readConfig() (config, error) {
	c := config{
		Config: server.DefaultConfig(),
	}
	if _, err := os.Stat("config.toml"); os.IsNotExist(err) {
		data, err := toml.Marshal(c)
		if err != nil {
			return c, fmt.Errorf("failed encoding default config: %v", err)
		}
		if err := os.WriteFile("config.toml", data, 0644); err != nil {
			return c, fmt.Errorf("failed creating config: %v", err)
		}
		return c, nil
	}
	data, err := os.ReadFile("config.toml")
	if err != nil {
		return c, fmt.Errorf("error reading config: %v", err)
	}
	if err := toml.Unmarshal(data, &c); err != nil {
		return c, fmt.Errorf("error decoding config: %v", err)
	}
	if c.SkyBlock.TreeType == "" {
		c.SkyBlock.TreeType = "oak"
	}
	return c, nil
}

type config struct {
	server.Config

	SkyBlock struct {
		// set tree type default is oak
		TreeType string `toml:"tree_type"`
	}
}
