package world

import (
	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
	"math"
	"math/rand"
)

// TreeInterface provides methods for Tree.
type TreeInterface interface {
	CanPlaceObject(w *world.World, x, y, z int, rand *rand.Rand) bool
	CanOverride(b world.Block) bool
	PlaceTree(w *world.World, pos *cube.Pos, rand *rand.Rand) error
	GenerateTrunkHeight(rand *rand.Rand) int
	PlaceTrunk(x, y, z int, rand *rand.Rand, trunkHeight int, w *world.World)
	PlaceCanopy(x, y, z int, rand *rand.Rand, w *world.World)
	GetName() string
}

// Tree is a tree that can be placed in the world.
type Tree struct {
	Name string
	// TrunkBlock is the block that the trunk of the tree is made of.
	TrunkBlock world.Block
	// LeafBlock is the block that the leaves of the tree are made of.
	LeafBlock world.Block
	// TreeHeight is the height of the tree.
	TreeHeight int
}

// BirchTree ...
type BirchTree struct {
	Tree

	SuperBirch bool
}

// SpruceTree ...
type SpruceTree struct {
	Tree
}

// CanPlaceObject determines if the tree can be placed at the given position.
func (t *Tree) CanPlaceObject(w *world.World, x, y, z int, rand *rand.Rand) bool {
	radiusToCheck := 0
	var xx, yy, zz int
	for yy = 0; yy < t.TreeHeight; y++ {
		if yy == 1 || yy == t.TreeHeight-1 {
			radiusToCheck++
		}
		for xx = -radiusToCheck; xx <= radiusToCheck; xx++ {
			for zz = -radiusToCheck; zz <= radiusToCheck; zz++ {
				if !t.CanOverride(w.Block(
					cube.Pos{
						x + xx,
						y + yy,
						z + zz,
					})) {
					return false
				}
			}
		}
	}
	return true
}

// CanOverride determines if the block can be overridden by the tree.
func (t *Tree) CanOverride(b world.Block) bool {
	_, ok := b.(*block.Wood)
	if ok {
		return true
	}
	_, ok = b.(*block.Leaves)
	if ok {
		return true
	}
	// TODO: Check block is solid
	return false
}

// PlaceTree places the tree at the given position.
func (t *Tree) PlaceTree(w *world.World, pos *cube.Pos, rand *rand.Rand) error {
	//if !t.CanPlaceObject(w, pos.X(), pos.Y(), pos.Z(), rand) {
	//	return fmt.Errorf("cannot place tree at %v", pos)
	//}
	t.PlaceTrunk(pos.X(), pos.Y(), pos.Z(), rand, t.GenerateTrunkHeight(rand), w)
	t.PlaceCanopy(pos.X(), pos.Y(), pos.Z(), rand, w)
	return nil
}

// GenerateTrunkHeight returns the height of the trunk of the tree.
func (t *Tree) GenerateTrunkHeight(rand *rand.Rand) int {
	return t.TreeHeight + 1
}

// PlaceTrunk places the trunk of the tree at the given position.
func (t *Tree) PlaceTrunk(x, y, z int, rand *rand.Rand, trunkHeight int, w *world.World) {
	w.SetBlock(cube.Pos{x, y - 1, z}, block.Dirt{}, nil)

	var yy int

	for yy = 0; yy < trunkHeight; yy++ {
		//if t.CanOverride(w.Block(cube.Pos{x, y + yy, z})) {
		w.SetBlock(cube.Pos{x, y + yy, z}, t.TrunkBlock, nil)
		//}
	}
}

// PlaceCanopy places the canopy of the tree at the given position.
func (t *Tree) PlaceCanopy(x, y, z int, rand *rand.Rand, w *world.World) {
	for yy := y - 3 + t.TreeHeight; yy <= y+t.TreeHeight; yy++ {
		yOff := yy - (y + t.TreeHeight)
		mid := 1 - yOff/2
		for xx := x - mid; xx <= x+mid; xx++ {
			xOff := int(math.Abs(float64(xx - x)))
			for zz := z - mid; zz <= z+mid; zz++ {
				zOff := int(math.Abs(float64(zz - z)))
				if xOff == mid && zOff == mid && (yOff == 0 || rand.Intn(2) == 0) {
					continue
				}
				//if !w.Block(cube.Pos{xx, yy, zz}).Solid() {
				w.SetBlock(cube.Pos{xx, yy, zz}, t.LeafBlock, nil)
				//}
			}
		}
	}
}

// PlaceCanopy places the canopy of the tree at the given position for SpruceTree.
func (t *SpruceTree) PlaceCanopy(x, y, z int, rand *rand.Rand, w *world.World) {
	topSize := t.TreeHeight - (1 + rand.Intn(2))
	lRadius := 2 + rand.Intn(2)
	radius := rand.Intn(2)
	maxR := 1
	minR := 0

	for yy := 0; yy <= topSize; yy++ {
		yyy := y + t.TreeHeight - yy
		for xx := x - radius; xx <= x+radius; xx++ {
			xOff := int(math.Abs(float64(xx - x)))
			for zz := z - radius; zz <= z+radius; zz++ {
				zOff := int(math.Abs(float64(zz - z)))
				if xOff == radius && zOff == radius && radius > 0 {
					continue
				}
				//if !w.Block(cube.Pos{xx, yyy, zz}).Solid() {
				w.SetBlock(cube.Pos{xx, yyy, zz}, t.LeafBlock, nil)
				//}
			}
		}
		if radius >= maxR {
			radius = minR
			minR = 1
			if maxR++; maxR > lRadius {
				maxR = lRadius
			}
		} else {
			radius++
		}
	}
}

// GenerateTrunkHeight returns height of the trunk of the SpruceTree.
func (t *SpruceTree) GenerateTrunkHeight(rand *rand.Rand) int {
	return t.TreeHeight + rand.Intn(3)
}

// NewOakTree returns a new oak Tree.
func NewOakTree() *Tree {
	return &Tree{
		Name: "Oak Tree",
		TrunkBlock: block.Wood{
			Wood:     block.OakWood(),
			Stripped: false,
			Axis:     cube.Y,
		},
		TreeHeight: 7,
		LeafBlock: block.Leaves{
			Wood:         block.OakWood(),
			Persistent:   true,
			ShouldUpdate: true,
		},
	}
}

// NewBirchTree returns a new BirchTree.
func NewBirchTree(superBirch bool) *BirchTree {
	return &BirchTree{
		SuperBirch: superBirch,
		Tree: Tree{
			Name: "Birch Tree",
			TrunkBlock: block.Wood{
				Wood:     block.BirchWood(),
				Stripped: false,
				Axis:     cube.Y,
			},
			TreeHeight: 7,
			LeafBlock: block.Leaves{
				Wood:         block.BirchWood(),
				Persistent:   true,
				ShouldUpdate: true,
			},
		},
	}
}

// NewSpruceTree returns a new SpruceTree.
func NewSpruceTree() *SpruceTree {
	return &SpruceTree{
		Tree: Tree{
			Name: "Spruce Tree",
			TrunkBlock: block.Wood{
				Wood:     block.SpruceWood(),
				Stripped: false,
				Axis:     cube.Y,
			},
			TreeHeight: 7,
			LeafBlock: block.Leaves{
				Wood:         block.SpruceWood(),
				Persistent:   true,
				ShouldUpdate: true,
			},
		},
	}
}

// NewJungleTree returns a new jungle Tree.
func NewJungleTree() *Tree {
	return &Tree{
		Name: "Jungle Tree",
		TrunkBlock: block.Wood{
			Wood:     block.JungleWood(),
			Stripped: false,
			Axis:     cube.Y,
		},
		LeafBlock: block.Leaves{
			Wood:         block.JungleWood(),
			Persistent:   false,
			ShouldUpdate: true,
		},
		TreeHeight: 7,
	}
}

func (t *Tree) GetName() string {
	return t.Name
}

func (t *BirchTree) GetName() string {
	return t.Name
}

func (t *SpruceTree) GetName() string {
	return t.Name
}
