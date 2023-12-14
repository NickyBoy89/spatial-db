package world

// `SectionPalette` is a "palette", which is a sort of look-up-table (LUT) between
// an index into the LUT, and the resulting `BlockID`
//
// This palette is unique to each section, and allows ranges of blocks to be
// changed in constant time, with the downside of having to "compact" the palette
type SectionPalette struct {
	ids map[PaletteIndex]BlockID
}

func NewSectionPalette() SectionPalette {
	return SectionPalette{
		ids: make(map[PaletteIndex]BlockID),
	}
}

// `Compact` removes all duplicate states from a palette and returns if the
// palette was modified
func (p *SectionPalette) Compact() bool {
	newIds := make(map[BlockID]PaletteIndex)

	var wasResized bool

	// Filter out the duplicate block ids
	for index, blockId := range p.ids {
		newIds[blockId] = index
	}

	// If there was not a resize, return instantly
	if len(newIds) != len(p.ids) {
		wasResized = true
	} else {
		return false
	}

	ids := make(map[PaletteIndex]BlockID)

	for blockId, index := range newIds {
		ids[index] = blockId
	}

	p.ids = ids

	return wasResized
}

// `IndexFor` returns the palette index for a specified block id
//
// If the block id does not exist, it is placed in the palette and the index
// is returned
func (p *SectionPalette) IndexFor(state BlockID) PaletteIndex {
	var maxIndex PaletteIndex

	// If the state is already in the palette, return it
	for index, blockId := range p.ids {
		if index > maxIndex {
			maxIndex = index
		}

		if state == blockId {
			return PaletteIndex(index)
		}
	}

	// Otherwise, insert it into the palette and return the index
	p.ids[maxIndex+1] = state

	return maxIndex + 1
}

// `State` returns the state for a palette's index
func (p SectionPalette) State(index PaletteIndex) BlockID {
	return p.ids[index]
}

// `ReplaceIndex` replaces the block state at the given palette' index
//
// This function is used to quickly zero a section with a given block state in
// constant time, and should not be used for any other purpose. Otherwise, set
// the block state within the section to the result of `IndexFor`.
func (p *SectionPalette) ReplaceIndex(index PaletteIndex, state BlockID) {
	p.ids[index] = state
}

// `Indexes` returns a list of all the palette indexes that are in the palette
func (p SectionPalette) Indexes() []PaletteIndex {
	inds := []PaletteIndex{}

	for ind := range p.ids {
		inds = append(inds, ind)
	}

	return inds
}

type PaletteIndex byte
