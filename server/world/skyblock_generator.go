package world

import (
	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/biome"
	"github.com/df-mc/dragonfly/server/world/chunk"
	"github.com/df-mc/dragonfly/server/world/generator"
)

type SkyBlockGenerator struct {
	generator.Flat
}

func NewGenerator() *SkyBlockGenerator {
	water := block.Water{
		Still:   true,
		Depth:   8,
		Falling: false,
	}
	return &SkyBlockGenerator{
		generator.NewFlat(biome.Ocean{}, []world.Block{water, water, water, block.Bedrock{}}),
	}
}

func (g *SkyBlockGenerator) GenerateChunk(pos world.ChunkPos, chunk *chunk.Chunk) {
	g.Flat.GenerateChunk(pos, chunk)
}
