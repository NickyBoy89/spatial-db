use crate::storage::world::{BlockID, BlockPos, BlockRange, ChunkData, ChunkPos};
use crate::storage_server::StorageServer;

#[derive(Debug)]
struct MultipleBlocks {
    id: BlockID,
    range: BlockRange,
}

#[derive(Debug)]
pub struct SimpleServer {
    chunks: Vec<ChunkData>,
    // block_ranges: Vec<MultipleBlocks>,
}

impl SimpleServer {
    pub fn new() -> Self {
        SimpleServer {
            chunks: Vec::new(),
            // block_ranges: Vec::new(),
        }
    }

    pub fn num_chunks(&self) -> usize {
        self.chunks.len()
    }

    fn chunk_at_block_mut(&mut self, block_pos: &BlockPos) -> Option<&mut ChunkData> {
        // Find what chunk the block is in
        let chunk_pos = ChunkPos::from(block_pos);

        // Find the chunk with the correct index
        for chunk in self.chunks.iter_mut() {
            if chunk.pos == chunk_pos {
                return Some(chunk);
            }
        }

        None
    }

    fn chunk_at(&self, block_pos: &BlockPos) -> Option<&ChunkData> {
        let chunk_pos = ChunkPos::from(block_pos);

        for chunk in self.chunks.iter() {
            if chunk.pos == chunk_pos {
                return Some(chunk);
            }
        }

        None
    }

    fn create_chunk_at(&mut self, chunk_pos: &ChunkPos) {
        let new_chunk = ChunkData::new(chunk_pos);

        self.chunks.push(new_chunk);
    }
}

impl StorageServer for SimpleServer {
    fn change_block(&mut self, target_state: BlockID, world_position: &BlockPos) {
        let mut chunk = self.chunk_at_block_mut(world_position);

        // Test if there is a chunk that already exists
        if chunk.is_none() {
            self.create_chunk_at(&ChunkPos::from(world_position));
            chunk = self.chunk_at_block_mut(world_position);
        }

        let chunk = chunk.expect("Could not find chunk");

        // Find the section that the block is located in
        let current_section = &mut chunk.sections[world_position.y % 16];
        // Find the index that the block is at, and update its state
        let chunk_array_index = current_section.index_of_block(&world_position);
        current_section.update_block_at_index(&target_state, chunk_array_index);
    }

    fn change_block_range(&mut self, target_stage: BlockID, start: &BlockPos, end: &BlockPos) {
        unimplemented!()
        // self.block_ranges.push(MultipleBlocks {
        //     id: target_stage,
        //     range: BlockRange {
        //         start: start.clone(),
        //         end: end.clone(),
        //     },
        // })
    }

    fn read_block_at(&self, pos: &BlockPos) -> BlockID {
        let chunk = self.chunk_at(pos);

        if let Some(chunk) = chunk {
            let chunk_section = chunk.section_for(pos);

            return chunk_section.get_block_at_index(pos).clone();
        }

        // for blocks in self.block_ranges.iter() {
        //     if blocks.range.within_range(&pos) {
        //         return blocks.id.clone();
        //     }
        // }

        BlockID::Empty
    }
}
