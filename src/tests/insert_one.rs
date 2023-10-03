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

        // Make sure the server only creates one chunk
        assert_eq!(server.num_chunks(), 1);

        // Retrieve one value
        assert_eq!(
            server.read_block_at(&BlockPos::new(0, 0, 0)),
            BlockID::Generic
        );

        // Retrieve an empty value in the current chunk
        assert_eq!(
            server.read_block_at(&BlockPos::new(1, 1, 1)),
            BlockID::Empty
        );

        // Retrieve a value in an empty chunk
        assert_eq!(
            server.read_block_at(&BlockPos::new(32, 32, 32)),
            BlockID::Empty
        );

        // Make sure a chunk was not created on that read
        assert_eq!(server.num_chunks(), 1);
    }

    #[test]
    fn test_remove_one_block() {
        let mut server = SimpleServer::new();

        let pos = BlockPos::new(0, 0, 0);

        server.change_block(BlockID::Generic, &pos);

        assert_eq!(server.num_chunks(), 1);

        assert_eq!(server.read_block_at(&pos), BlockID::Generic);

        server.change_block(BlockID::Empty, &pos);

        assert_eq!(server.read_block_at(&pos), BlockID::Empty);
    }

    #[test]
    fn test_insert_some_blocks() {
        let mut server = SimpleServer::new();

        let blocks = [
            BlockPos::new(0, 2, 0),
            BlockPos::new(0, 2, 1),
            BlockPos::new(0, 2, -1),
            BlockPos::new(1, 2, 0),
            BlockPos::new(-1, 2, 0),
            BlockPos::new(0, 3, 0),
            BlockPos::new(0, 0, 0),
        ];

        for pos in blocks.iter() {
            server.change_block(BlockID::Generic, pos);
        }

        assert_eq!(server.num_chunks(), 1);

        for pos in blocks.iter() {
            let read = server.read_block_at(pos);
            println!("Pos: {:?}, {:?}", pos, read);
            assert_eq!(read, BlockID::Generic);
        }
    }
}
