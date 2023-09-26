type ChunkCoordinate = isize;

const SECTIONS_PER_CHUNK: usize = 16;

struct ChunkData {
    x: ChunkCoordinate,
    y: ChunkCoordinate,
    sections: [ChunkSection; SECTIONS_PER_CHUNK],
}

// https://wiki.vg/Chunk_Format
struct ChunkSection {
    /// The number of non-empty blocks in the section. If completely full, the
    /// section contains a 16 x 16 x 16 cube of blocks = 4096 blocks
    /// If the section is empty, this is skipped
    block_count: u16,
    /// The data for all the blocks in the chunk
    /// The representation for this may be different based on the number of
    /// non-empty blocks
    block_states: [BlockID; 4096],
}

/// `BlockPos` represents the location of a block in world space
pub struct BlockPos {
    x: isize,
    y: isize,
    z: isize,
}

/// BlockID represents the type of block stored
#[repr(u8)]
pub enum BlockID {
    Empty,
    Generic,
}
