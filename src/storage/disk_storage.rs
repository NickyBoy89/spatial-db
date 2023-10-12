use super::world::{BlockID, ChunkData, ChunkPos};
use std::fs::File;

const DATABASE_FILE_LOCATION: &str = "./persistence";

struct RunLengthEncoding {
    pairs: Vec<(usize, BlockID)>,
}

impl RunLengthEncoding {
    fn from_chunk(chunk_data: &ChunkData) -> Self {
        for section in chunk_data.sections {
            for index in section.chunk_data {
                // Yes
            }
        }
    }
}

