use crate::simple_server::server::SimpleServer;
use crate::storage::world::{BlockPos, BlockRange};

#[cfg(test)]
mod tests {
    use crate::{storage::world::BlockID, storage_server::StorageServer};

    use super::*;
    #[test]
    fn within_two_dimensions() {
        // Get two points on the same z axis
        let first = BlockPos::new(0, 0, 0);
        let second = BlockPos::new(4, 4, 0);

        let range = BlockRange::new(&first, &second);

        let test1 = BlockPos::new(1, 1, 0);
        let test2 = BlockPos::new(0, 0, 0);
        let test3 = BlockPos::new(0, 4, 0);
        let test4 = BlockPos::new(4, 4, 0);
        let test5 = BlockPos::new(4, 0, 0);

        assert!(range.within_range(&test1));
        assert!(range.within_range(&test2));
        assert!(range.within_range(&test3));
        assert!(range.within_range(&test4));
        assert!(range.within_range(&test5));

        let test6 = BlockPos::new(-1, 0, 0);

        assert!(!range.within_range(&test6));
    }

    #[test]
    fn test_simple_insert() {
        let mut server = SimpleServer::new();

        server.change_block(BlockID::Generic, &BlockPos::new(0, 0, 0));

        assert_eq!(
            server.read_block_at(&BlockPos::new(0, 0, 0)),
            BlockID::Generic
        );
    }
}
