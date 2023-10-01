use crate::storage::world::{BlockID, BlockPos, BlockRange, ChunkData, ChunkPos};
use crate::storage_server::StorageServer;

#[derive(Clone)]
struct SingleBlock {
    id: BlockID,
    position: BlockPos,
}

struct MultipleBlocks {
    id: BlockID,
    range: BlockRange,
}

pub struct SimpleServer {
    chunks: Vec<ChunkData>,
    single_blocks: Vec<SingleBlock>,
    block_ranges: Vec<MultipleBlocks>,
}

impl SimpleServer {
    pub fn new() -> Self {
        SimpleServer {
            chunks: Vec::new(),
            single_blocks: Vec::new(),
            block_ranges: Vec::new(),
        }
    }
}

impl StorageServer for SimpleServer {
    fn change_block(&mut self, target_state: BlockID, world_position: BlockPos) {
        let chunk_pos = ChunkPos::from(&world_position);

        println!("Chunk position: {:?}", chunk_pos);

        for chunk in self.chunks.iter_mut() {
            if chunk.pos.same_location(&chunk_pos) {
                let current_section = &mut chunk.sections[world_position.y / 16];
                let chunk_array_index = current_section.index_of_block(&world_position);
                current_section.update_block_at_index(&target_state, chunk_array_index);
            }
        }

        self.single_blocks.push(SingleBlock {
            id: target_state,
            position: world_position,
        });
    }

    fn change_block_range(&mut self, target_stage: BlockID, start: BlockPos, end: BlockPos) {
        self.block_ranges.push(MultipleBlocks {
            id: target_stage,
            range: BlockRange { start, end },
        })
    }

    fn read_block_at(&self, pos: BlockPos) -> BlockID {
        for block in self.single_blocks.iter() {
            if block.position == pos {
                return block.id.clone();
            }
        }

        for blocks in self.block_ranges.iter() {
            if blocks.range.within_range(&pos) {
                return blocks.id.clone();
            }
        }

        BlockID::Empty
    }
}
