package world

import (
	"github.com/Tnze/go-mc/save"
)

const (
	// The number of sections per chunk. This determines the total height of the
	// chunk
	ChunkSectionCount = 16

	// The number of blocks in a horizontal slice of a chunk
	chunkSliceSize = 16 * 16
)

// `ChunkData` represents the contents of a "chunk", which is a column of voxels
// in world space
type ChunkData struct {
	// The position of the chunk, in world space
	Pos ChunkPos `json:"pos"`
	// The column of sections
	Sections [ChunkSectionCount]ChunkSection `json:"sections"`
}

// `ChunkSection' is a fixed-size cube that stores the data in a chunk
type ChunkSection struct {
	// A look-up-table of each section index to its value
	Palette     SectionPalette             `json:"palette"`
	BlockStates [16 * 16 * 16]PaletteIndex `json:"block_states"`
}

func rem_euclid(a, b int) int {
	return (a%b + b) % b
}

func IndexOfBlock(pos BlockPos) int {
	baseX := rem_euclid(pos.X, 16)
	baseY := rem_euclid(int(pos.Y), 16)
	baseZ := rem_euclid(pos.Z, 16)

	return (baseY * chunkSliceSize) + (baseZ * 16) + baseX
}

func (cs *ChunkSection) UpdateBlock(pos BlockPos, targetState BlockID) {
	cs.BlockStates[IndexOfBlock(pos)] = cs.Palette.IndexFor(targetState)
}

func (cs *ChunkSection) FetchBlock(pos BlockPos) BlockID {
	return cs.Palette.State(cs.BlockStates[IndexOfBlock(pos)])
}

func (cd *ChunkData) SectionFor(pos BlockPos) *ChunkSection {
	return &cd.Sections[pos.Y%ChunkSectionCount]
}

func (cd *ChunkData) IndexToBlockPos(index int) BlockPos {
	posX := index % 16
	posZ := ((index - posX) % 256) / 16
	posY := ((index - posZ) % 4096) / 256
	return BlockPos{
		X: posX + (cd.Pos.X * 16),
		Y: uint(posY),
		Z: posZ + (cd.Pos.Z * 16),
	}
}

// Conversion from Minecraft chunks

func extractPaletteIndexes(compressed int64) [16]byte {
	var outputs [16]byte
	var outputIndex int

	for index := 0; index < 64; index += 4 {
		shifted := compressed >> index
		// Mask off the lowest four bits
		shifted &= 0xf
		outputs[outputIndex] = byte(shifted)
		outputIndex += 1
	}

	return outputs
}

func (cd *ChunkData) FromMCAChunk(other save.Chunk) {
	// Load the chunk's position
	cd.Pos = ChunkPos{X: int(other.XPos), Z: int(other.ZPos)}

	// Load the data from the chunk
	for sectionIndex, section := range other.Sections {
		// TODO: Enable chunks to have more than 16 sections
		if sectionIndex >= ChunkSectionCount {
			break
		}

		var currentSection ChunkSection
		currentSection.Palette = NewSectionPalette()

		paletteIndexes := []int{}
		for _, compress := range section.BlockStates.Data {
			indexes := extractPaletteIndexes(int64(compress))

			converted := make([]int, 16)
			for i := 0; i < 16; i++ {
				converted[i] = int(indexes[i])
			}

			paletteIndexes = append(paletteIndexes, converted...)
		}

		for blockIndex, paletteIndex := range paletteIndexes {
			var state BlockID
			if section.BlockStates.Palette[paletteIndex].Name == "minecraft:air" {
				state = Empty
			} else {
				state = Generic
			}

			// TODO: Remove this workaround for larger bit sizes in palettes
			if blockIndex < 4096 {
				currentSection.BlockStates[blockIndex] = currentSection.Palette.IndexFor(state)
			}
		}

		cd.Sections[sectionIndex] = currentSection
	}
}
