use core::fmt;
use serde::ser;
use serde::ser::{SerializeSeq, SerializeStruct, Serializer};
use serde::Serialize;
use std::{
    cmp::{max, min},
    fmt::Debug,
    fs::File,
};

const SECTIONS_PER_CHUNK: usize = 16;
const SLICE_SIZE: usize = 16 * 16;

#[derive(Debug, Clone, PartialEq, Serialize)]
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
    pub fn storage_file_name(&self) -> String {
        format!("{}.{}.chunk", self.x, self.z)
    }
}

#[derive(Debug)]
pub struct ChunkData {
    pub pos: ChunkPos,
    pub sections: [ChunkSection; SECTIONS_PER_CHUNK],
}

impl Serialize for ChunkData {
    fn serialize<S>(&self, serializer: S) -> Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        let mut seq = serializer.serialize_seq(Some(self.sections.len()))?;

        for section in self.sections {
            seq.serialize_element(&section)?;
        }
        seq.end()
    }
}

impl ChunkData {
    pub fn new(pos: &ChunkPos) -> Self {
        ChunkData {
            pos: pos.clone(),
            sections: [ChunkSection::new(); SECTIONS_PER_CHUNK],
        }
    }

    pub fn section_for(&self, block_pos: &BlockPos) -> &ChunkSection {
        &self.sections[block_pos.y % 16]
    }

    pub fn write_to_file(&self, output_file: &mut File) {
        let serialized = serde_json::to_string(self).unwrap();
    }

    pub fn read_from_file(chunk_file: &File) -> Self {
        unimplemented!()
    }
}

// https://wiki.vg/Chunk_Format
#[derive(Clone, Copy, Serialize)]
pub struct ChunkSection {
    /// The number of non-empty blocks in the section. If completely full, the
    /// section contains a 16 x 16 x 16 cube of blocks = 4096 blocks
    /// If the section is empty, this is skipped
    block_count: u16,
    /// The data for all the blocks in the chunk
    /// The representation for this may be different based on the number of
    /// non-empty blocks
    #[serde(with = "serde_arrays")]
    block_states: [BlockID; 16 * 16 * 16],
}

impl Debug for ChunkSection {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        write!(f, "ChunkSection {{ blocks: {}, states: ", self.block_count)?;
        if self.block_count > 0 {
            write!(f, "{:?}", self.block_states)?;
        }
        write!(f, " }}")
    }
}

impl ChunkSection {
    pub fn new() -> Self {
        ChunkSection {
            block_count: 0,
            block_states: [BlockID::Empty; 16 * 16 * 16],
        }
    }

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

    pub fn get_block_at_index(&self, pos: &BlockPos) -> &BlockID {
        let array_index = self.index_of_block(pos);

        &self.block_states[array_index]
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
#[derive(Debug)]
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
#[derive(Debug, Clone, Copy, PartialEq, Serialize)]
pub enum BlockID {
    Empty,
    Generic,
}
