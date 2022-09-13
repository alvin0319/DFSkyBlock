package world

import (
	"fmt"
	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/mcdb"
	"github.com/df-mc/goleveldb/leveldb/opt"
	"github.com/sirupsen/logrus"
	"math/rand"
	"time"
)

type Manager struct {
	// worlds holds all worlds that are currently loaded.
	worlds map[string]*world.World

	Log *logrus.Logger

	Tree TreeInterface
}

// NewManager returns a new world manager.
func NewManager(treeType TreeInterface) *Manager {
	return &Manager{
		worlds: make(map[string]*world.World),
		Log:    logrus.New(),
		Tree:   treeType,
	}
}

// GetWorld returns the world by the given name. If the world is not loaded, it will return error.
func (m *Manager) GetWorld(name string) (*world.World, error) {
	if w, ok := m.worlds[name]; ok {
		return w, nil
	}
	return nil, fmt.Errorf("world %v does not exist or has not been loaded yet", name)
}

// NewWorld generates new world by given name.
// If creating provider fails, the error will be returned.
func (m *Manager) NewWorld(name string) (*world.World, error) {
	p, err := mcdb.New("skyblock/"+name, opt.DefaultCompression)
	if err != nil {
		return nil, err
	}
	conf := world.Config{
		Log: m.Log,
		Dim: world.Overworld,
		PortalDestination: func(dim world.Dimension) *world.World {
			return nil
		},
		Provider:        p,
		Generator:       NewGenerator(),
		ReadOnly:        false,
		RandomTickSpeed: 3,
		RandSource:      rand.NewSource(time.Now().Unix()),
	}
	w := conf.New()

	yStart := -63

	for x := 0; x < 16; x++ {
		for z := 0; z < 16; z++ {
			for y := yStart; y < yStart+5; y++ {
				if y < yStart+4 {
					w.SetBlock(cube.Pos{x, y, z}, block.Dirt{}, nil)
				} else {
					w.SetBlock(cube.Pos{x, y, z}, block.Grass{}, nil)
				}
			}
		}
	}

	w.SetSpawn(cube.Pos{3, yStart + 4, 8})

	treePos := cube.Pos{8, yStart + 4, 8}
	err = m.Tree.PlaceTree(w, &treePos, rand.New(rand.NewSource(time.Now().Unix())))
	if err != nil {
		return nil, err
	}

	return w, nil
}
