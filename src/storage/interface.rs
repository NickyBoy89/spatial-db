use crate::storage::world::{BlockID, BlockPos};

enum StorageInterface {
    /// `ChangeBlock` changes the block at the world position given by `world_position` to the
    /// target block id `BlockID`
    ChangeBlock {
        target_state: BlockID,
        world_position: BlockPos,
    },
    ChangeBlockRange(BlockID, BlockPos, BlockPos),
}
