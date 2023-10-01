use std::cmp::{max, min};

const SECTIONS_PER_CHUNK: usize = 16;
const SLICE_SIZE: usize = 16 * 16;

#[derive(Debug)]
pub struct ChunkPos {
    pub x: isize,
    pub z: isize,
}

impl From<&BlockPos> for ChunkPos {
    fn from(value: &BlockPos) -> Self {
        ChunkPos {
            x: value.x / 16,
            z: value.z / 16,
        }
    }
}

impl ChunkPos {
    pub fn same_location(&self, other: &ChunkPos) -> bool {
        self.x == other.x && self.z == other.z
    }
}

pub struct ChunkData {
    pub pos: ChunkPos,
    pub sections: [ChunkSection; SECTIONS_PER_CHUNK],
}

// https://wiki.vg/Chunk_Format
pub struct ChunkSection {
    /// The number of non-empty blocks in the section. If completely full, the
    /// section contains a 16 x 16 x 16 cube of blocks = 4096 blocks
    /// If the section is empty, this is skipped
    block_count: u16,
    /// The data for all the blocks in the chunk
    /// The representation for this may be different based on the number of
    /// non-empty blocks
    block_states: [BlockID; 16 * 16 * 16],
}

impl ChunkSection {
    pub fn index_of_block(&self, pos: &BlockPos) -> usize {
        let base_x = pos.x.rem_euclid(16) as usize;
        let base_y = pos.y.rem_euclid(16) as usize;
        let base_z = pos.z.rem_euclid(16) as usize;

        (base_y * SLICE_SIZE) + (base_z * 16) + base_x
    }

    pub fn update_block_at_index(&mut self, id: &BlockID, index: usize) {
        let existing_block = &self.block_states[index];
        match existing_block {
            BlockID::Empty => match id {
                BlockID::Generic => {
                    // If the existing block is empty, and the block that we
                    // are inserting is non-empty, increment the number of blocks
                    self.block_count += 1;
                }
                _ => {}
            },
            _ => match id {
                BlockID::Empty => {
                    // If the existing block is non-empty, and the block that
                    // we are inserting is empty, then decrement the number of
                    // blocks
                    self.block_count -= 1;
                }
                _ => {}
            },
        }

        self.block_states[index] = id.clone();
    }
}

/// `BlockPos` represents the location of a block in world space
#[derive(Debug, Clone, PartialEq)]
pub struct BlockPos {
    pub x: isize,
    pub y: usize,
    pub z: isize,
}

impl BlockPos {
    pub fn new(x: isize, y: usize, z: isize) -> Self {
        BlockPos { x, y, z }
    }
}

/// BlockRange represents a range of blocks that have been updated
pub struct BlockRange {
    pub start: BlockPos,
    pub end: BlockPos,
}

impl BlockRange {
    pub fn new(start: &BlockPos, end: &BlockPos) -> Self {
        BlockRange {
            start: start.clone(),
            end: end.clone(),
        }
    }
    pub fn within_range(&self, pos: &BlockPos) -> bool {
        let minx = min(self.start.x, self.end.x);
        let maxx = max(self.start.x, self.end.x);

        if pos.x < minx || pos.x > maxx {
            return false;
        }

        let miny = min(self.start.y, self.end.y);
        let maxy = max(self.start.y, self.end.y);

        if pos.y < miny || pos.y > maxy {
            return false;
        }

        let minz = min(self.start.z, self.end.z);
        let maxz = max(self.start.z, self.end.z);

        if pos.z < minz || pos.z > maxz {
            return false;
        }

        true
    }
}

/// BlockID represents the type of block stored
#[repr(u8)]
#[derive(Debug, Clone, PartialEq)]
pub enum BlockID {
    Empty,
    Generic,
}
