use super::world::{BlockID, ChunkData, ChunkPos};
use std::cmp::Ordering;
use std::{collections::HashMap, fs::File, time::Instant};

const DATABASE_FILE_LOCATION: &str = "./persistence";

struct ChunkFile {}

const CACHED_CHUNK_FILES: usize = 1;

/// `ChunkStorageCache` caches a list of the most recently used file handles
/// where chunks are stored from, and allows for faster accessing of the data
/// from chunks
struct ChunkStorageCache {
    // `cached_chunk_files` is a vector of cached file handles that are already open
    cached_chunk_files: [File; CACHED_CHUNK_FILES],
    // `cached_file_names` is a list of all the filenames that are contained
    // within the cache
    cached_file_names: HashMap<String, usize>,
    last_used_times: [Instant; CACHED_CHUNK_FILES],
}

impl ChunkStorageCache {
    fn load_chunk_file(&mut self, file_name: &str) -> &File {
        let chunk_file = File::open(file_name).expect("Opening file for chunk failed");

        // Add the newly opened file to the cache

        // Insert the new item to replace the item that was last accessed
        // The minimum time should be the oldest time
        let (last_used_index, _) = self
            .last_used_times
            .iter()
            .enumerate()
            .reduce(
                |(fst_index, fst_time), (snd_index, snd_time)| match fst_time.cmp(&snd_time) {
                    Ordering::Less => (fst_index, fst_time),
                    Ordering::Equal | Ordering::Greater => (snd_index, snd_time),
                },
            )
            .expect("There should always be a last used index");

        // Next, we have to:
        // * Remove the old filename and index mapping from the names
        // * Replace the last used time with the curent time
        // * Replace the open file with the current one

        // Find the name of the previous entry
        let (previous_file_name, _) = self
            .cached_file_names
            .iter()
            .find(|(_, &array_index)| array_index == last_used_index)
            .expect("The last used index should always have a name");

        self.cached_file_names.remove(&previous_file_name.clone());
        self.cached_file_names
            .insert(file_name.to_string(), last_used_index);
        // Replace the timestamp with the new timestamp
        self.last_used_times[last_used_index] = Instant::now();
        self.cached_chunk_files[last_used_index] = chunk_file;

        &self.cached_chunk_files[last_used_index]
    }

    /// `fetch_chunk_by_pos` takes in the position of a chunk, and returns the
    /// data of the chunk from disk
    ///
    /// This operation is cached, if possible, so that subsequent accesses to
    /// the same chunk are handled by the same file
    pub fn fetch_chunk_by_pos(&mut self, pos: &ChunkPos) -> ChunkData {
        let file_name = pos.storage_file_name();

        let file_index = self.cached_file_names.get(file_name.as_str());

        let chunk_file = match file_index {
            Some(index) => &self.cached_chunk_files[*index],
            None => self.load_chunk_file(file_name.as_str()),
        };

        panic!("Yes");
    }
}
