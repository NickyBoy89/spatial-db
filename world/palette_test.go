package world

import (
	"runtime/debug"
	"testing"
)

func newPaletteWith(ids []BlockID) *SectionPalette {
	p := NewSectionPalette()

	for _, id := range ids {
		p.IndexFor(id)
	}

	return &p
}

func checkPaletteEqual(t *testing.T, p *SectionPalette, index PaletteIndex, expected BlockID) {
	if p.State(index) != expected {
		debug.PrintStack()
		t.Fatalf("Expected to get state %v at index %v, got %v", expected, index, p.State(index))
	}
}

func TestInsertPalette(t *testing.T) {
	p := NewSectionPalette()

	ids := []BlockID{Empty, Generic}

	for _, id := range ids {
		index := p.IndexFor(id)
		if id != p.State(index) {
			t.Fatalf("Fetching index for id %v: got index %v which returned %v", id, index, p.State(index))
		}
	}
}

func TestReplaceCompactPalettte(t *testing.T) {
	p := newPaletteWith([]BlockID{Empty, Generic})

	zeroIndex := p.IndexFor(Empty)
	checkPaletteEqual(t, p, zeroIndex, Empty)

	// Zero out the chunk
	genIndex := p.IndexFor(Generic)
	checkPaletteEqual(t, p, genIndex, Generic)
	p.ReplaceIndex(genIndex, Empty)

	// Check that everything now returns zero
	checkPaletteEqual(t, p, zeroIndex, Empty)
	checkPaletteEqual(t, p, genIndex, Empty)

	// Test for compaction
	wasCompacted := p.Compact()
	if !wasCompacted {
		t.Fatalf("Palette should have been compacted")
	}

	// Test to see that we don't compact again
	wasCompacted = p.Compact()
	if wasCompacted {
		t.Fatalf("Palette should have not been compacted again")
	}
}
