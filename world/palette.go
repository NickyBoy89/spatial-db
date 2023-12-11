package world

// `SectionPalette` is a "palette", which is a sort of look-up-table (LUT) between
// an index into the LUT, and the resulting `BlockID`
//
// This palette is unique to each section, and allows ranges of blocks to be
// changed in constant time, with the downside of having to "compact" the palette
type SectionPalette []BlockID

// `Compact` removes all duplicate states from a palette and returns a new palette
func (p SectionPalette) Compact() SectionPalette {
	ids := make(map[BlockID]bool)

	// Filter out the duplicate block ids
	for _, blockId := range p {
		ids[blockId] = true
	}

	var np SectionPalette

	for blockId := range ids {
		np = append(np, blockId)
	}

	return np
}

// `IndexFor` returns the palette index for a specified block id
//
// If the block id does not exist, it is placed in the palette and the index
// is returned
func (p SectionPalette) IndexFor(state BlockID) PaletteIndex {
	// If the state is already in the palette, return it
	for index, blockId := range p {
		if state == blockId {
			return PaletteIndex(index)
		}
	}

	// Otherwise, insert it into the palette and return the index
	p = append(p, state)

	return PaletteIndex(len(p) - 1)
}

// `State` returns the state for a palette's index
func (p SectionPalette) State(index PaletteIndex) BlockID {
	return p[index]
}

type PaletteIndex byte
