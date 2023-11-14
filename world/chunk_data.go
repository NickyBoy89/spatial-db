package world

import (
	"github.com/Tnze/go-mc/save"
)

const ChunkSectionCount = 16

type ChunkData struct {
	Pos      ChunkPos                        `json:"pos"`
	Sections [ChunkSectionCount]ChunkSection `json:"sections"`
}

func (cd *ChunkData) SectionFor(pos BlockPos) *ChunkSection {
	return &cd.Sections[pos.Y%ChunkSectionCount]
}

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
				cd.Sections[sectionIndex].BlockCount += 1
				state = Generic
			}

			// TODO: Remove this workaround for larger bit sizes in palettes
			if blockIndex < 4096 {
				currentSection.BlockStates[blockIndex] = state
			}
		}

		cd.Sections[sectionIndex] = currentSection
	}
}
