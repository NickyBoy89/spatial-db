use crate::storage::world::{BlockID, BlockPos};
use axum::{routing::get, Router};

pub trait StorageServer {
    /// `change_block` changes the block at the world position given by `world_position` to the
    /// target block id `BlockID`
    fn change_block(&mut self, target_state: BlockID, world_position: &BlockPos);
    fn change_block_range(&mut self, target_stage: BlockID, start: &BlockPos, end: &BlockPos);

    /// `read_block_at` returns the id of the block at the location specified
    /// If no block is present, the returned id will be of the empty type
    fn read_block_at(&self, pos: &BlockPos) -> BlockID;
}

#[tokio::main]
pub async fn main() {
    let app = Router::new().route("/", get(|| async { "Hello World" }));

    axum::Server::bind(&"0.0.0.0:5000".parse().unwrap())
        .serve(app.into_make_service())
        .await
        .unwrap();
}
