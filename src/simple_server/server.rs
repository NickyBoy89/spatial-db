use crate::storage::disk_storage::ChunkStorageCache;
use crate::storage::world::{BlockID, BlockPos, BlockRange, ChunkData, ChunkPos};
use crate::storage_server::StorageServer;

#[derive(Debug)]
struct MultipleBlocks {
    id: BlockID,
    range: BlockRange,
}

#[derive(Debug)]
pub struct SimpleServer {
    chunk_storage: ChunkStorageCache,
}

impl SimpleServer {
    pub fn new() -> Self {
        SimpleServer {
            chunk_storage: ChunkStorageCache::new(),
        }
    }

    pub fn num_chunks(&self) -> usize {
        unimplemented!()
    }

    fn chunk_at(&mut self, block_pos: &BlockPos) -> Option<ChunkData> {
        let chunk_pos = ChunkPos::from(block_pos);

        let chunk = self
            .chunk_storage
            .fetch_chunk_by_pos(&chunk_pos)
            .expect("Finding chunk failed");

        Some(chunk)
    }

    fn create_chunk_at(&mut self, chunk_pos: &ChunkPos) {
        self.chunk_storage
            .fetch_chunk_by_pos(&chunk_pos)
            .expect("Creatinc chunk failed");
    }
}

impl StorageServer for SimpleServer {
    fn change_block(&mut self, target_state: BlockID, world_position: &BlockPos) {
        let mut chunk = self.chunk_at(world_position);

        // Test if there is a chunk that already exists
        if chunk.is_none() {
            self.create_chunk_at(&ChunkPos::from(world_position));
            chunk = self.chunk_at(world_position);
        }

        let mut chunk = chunk.expect("Could not find chunk");

        // Find the section that the block is located in
        let current_section = &mut chunk.sections[world_position.y % 16];
        // Find the index that the block is at, and update its state
        let chunk_array_index = current_section.index_of_block(&world_position);
        current_section.update_block_at_index(&target_state, chunk_array_index);
    }

    fn change_block_range(&mut self, target_stage: BlockID, start: &BlockPos, end: &BlockPos) {
        unimplemented!()
    }

    fn read_block_at(&mut self, pos: &BlockPos) -> BlockID {
        let chunk = self.chunk_at(pos);

        if let Some(chunk) = chunk {
            let chunk_section = chunk.section_for(pos);

            return chunk_section.get_block_at_index(pos).clone();
        }

        BlockID::Empty
    }
}
