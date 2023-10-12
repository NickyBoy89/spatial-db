extern crate test;

#[cfg(test)]
mod tests {
    use super::*;
    use crate::{
        simple_server::server::SimpleServer,
        storage::world::{BlockID, BlockPos},
        storage_server::StorageServer,
    };
    use rand::prelude::*;
    use test::Bencher;

    #[bench]
    fn bench_add_sequential_elements(b: &mut Bencher) {
        let mut server = SimpleServer::new();
        let mut x = 0;

        b.iter(|| {
            server.change_block(BlockID::Generic, &BlockPos::new(x, 0, 0));
            x += 1;
        });
    }

    #[bench]
    fn bench_add_clustered_points(b: &mut Bencher) {
        let mut server = SimpleServer::new();

        let mut rng = rand::thread_rng();

        static MAX_RANGE: isize = 128;

        b.iter(|| {
            let x: isize = rng.gen_range(-MAX_RANGE..MAX_RANGE);
            let y: usize = rng.gen::<u8>() as usize;
            let z: isize = rng.gen_range(-MAX_RANGE..MAX_RANGE);

            server.change_block(BlockID::Generic, &BlockPos::new(x, y, z));
        });
    }

    #[bench]
    fn bench_add_spread_out_points(b: &mut Bencher) {
        let mut server = SimpleServer::new();

        let mut rng = rand::thread_rng();

        static MAX_RANGE: isize = 65536;

        b.iter(|| {
            let x: isize = rng.gen_range(-MAX_RANGE..MAX_RANGE);
            let y: usize = rng.gen::<u8>() as usize;
            let z: isize = rng.gen_range(-MAX_RANGE..MAX_RANGE);

            server.change_block(BlockID::Generic, &BlockPos::new(x, y, z));
        });
    }

    #[bench]
    fn bench_insert_and_read_clustered(b: &mut Bencher) {
        let mut server = SimpleServer::new();
        let mut rng = rand::thread_rng();

        static NUM_BLOCKS: usize = 1_000;
        static MAX_RANGE: isize = 128;

        let mut positions = Vec::with_capacity(NUM_BLOCKS);

        for _ in 0..NUM_BLOCKS {
            let x: isize = rng.gen_range(-MAX_RANGE..MAX_RANGE);
            let y: usize = rng.gen::<u8>() as usize;
            let z: isize = rng.gen_range(-MAX_RANGE..MAX_RANGE);

            let pos = BlockPos::new(x, y, z);

            server.change_block(BlockID::Generic, &BlockPos::new(x, y, z));

            positions.push(pos);
        }

        b.iter(|| {
            for i in 0..NUM_BLOCKS {
                assert_eq!(server.read_block_at(&positions[i]), BlockID::Generic);
            }
        });
    }

    #[bench]
    fn bench_insert_and_read_cache(b: &mut Bencher) {
        let mut server = SimpleServer::new();
        let mut rng = rand::thread_rng();

        static NUM_BLOCKS: usize = 1_000;
        static MAX_RANGE: isize = 128;
        static EXPANDED_RANGE: isize = 2048;

        let mut positions = Vec::with_capacity(NUM_BLOCKS);

        for _ in 0..NUM_BLOCKS {
            let x: isize = rng.gen_range(-MAX_RANGE..MAX_RANGE);
            let y: usize = rng.gen::<u8>() as usize;
            let z: isize = rng.gen_range(-MAX_RANGE..MAX_RANGE);

            let pos = BlockPos::new(x, y, z);

            server.change_block(BlockID::Generic, &BlockPos::new(x, y, z));

            positions.push(pos);
        }

        b.iter(|| {
            // Read blocks that are already in the server
            for i in 0..NUM_BLOCKS {
                assert_eq!(server.read_block_at(&positions[i]), BlockID::Generic);
            }

            // Read blocks that might not be in the server, triggering a miss
            for _ in 0..NUM_BLOCKS {
                let x: isize = rng.gen_range(-EXPANDED_RANGE..EXPANDED_RANGE);
                let y: usize = rng.gen::<u8>() as usize;
                let z: isize = rng.gen_range(-EXPANDED_RANGE..EXPANDED_RANGE);
                server.read_block_at(&BlockPos::new(x, y, z));
            }
        });
    }

    #[bench]
    fn bench_clustered_many_misses(b: &mut Bencher) {
        let mut server = SimpleServer::new();
        let mut rng = rand::thread_rng();

        static NUM_BLOCKS: usize = 1_000;
        static MAX_RANGE: isize = 128;
        static EXPANDED_RANGE: isize = 2048;

        for _ in 0..NUM_BLOCKS {
            let x: isize = rng.gen_range(-MAX_RANGE..MAX_RANGE);
            let y: usize = rng.gen::<u8>() as usize;
            let z: isize = rng.gen_range(-MAX_RANGE..MAX_RANGE);

            server.change_block(BlockID::Generic, &BlockPos::new(x, y, z));
        }

        b.iter(|| {
            // Read blocks that might not be in the server, triggering a miss
            for _ in 0..NUM_BLOCKS {
                let x: isize = rng.gen_range(-EXPANDED_RANGE..EXPANDED_RANGE);
                let y: usize = rng.gen::<u8>() as usize;
                let z: isize = rng.gen_range(-EXPANDED_RANGE..EXPANDED_RANGE);
                server.read_block_at(&BlockPos::new(x, y, z));
            }
        });
    }
}
