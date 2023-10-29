use super::world::{ChunkData, ChunkPos};
use std::cmp::Ordering;
use std::error::Error;
use std::io::{BufReader, ErrorKind, Write};
use std::{collections::HashMap, fs::File, time::Instant};

const CACHED_CHUNK_FILES: usize = 1;

/// `ChunkStorageCache` caches a list of the most recently used file handles
/// where chunks are stored from, and allows for faster accessing of the data
/// from chunks
#[derive(Debug)]
pub struct ChunkStorageCache {
    // `cached_chunk_files` is a vector of cached file handles that are already open
    cached_chunk_files: [Option<File>; CACHED_CHUNK_FILES],
    // `cached_file_names` is a list of all the filenames that are contained
    // within the cache
    cached_file_names: HashMap<String, usize>,
    last_used_times: [Instant; CACHED_CHUNK_FILES],
}

impl ChunkStorageCache {
    pub fn new() -> Self {
        ChunkStorageCache {
            cached_chunk_files: [None; CACHED_CHUNK_FILES],
            cached_file_names: HashMap::new(),
            last_used_times: [Instant::now(); CACHED_CHUNK_FILES],
        }
    }

    /// `load_chunk_file` is called whenever a file is missing in the file cache
    /// and needs to be loaded from disk
    ///
    /// This replaces a slot for another file in the cache, according to the
    /// caching strategy
    fn load_chunk_file(&mut self, chunk_pos: &ChunkPos, file_name: &str) -> &File {
        let chunk_file = File::options().write(true).read(true).open(file_name);

        let chunk_file = match chunk_file {
            Ok(file) => file,
            Err(err) => match err.kind() {
                ErrorKind::NotFound => {
                    let mut new_chunk_file = File::options()
                        .write(true)
                        .read(true)
                        .create(true)
                        .open(file_name)
                        .expect("Opening new chunk file failed");

                    let blank_chunk = ChunkData::new(chunk_pos);

                    let encoded_chunk = serde_json::to_string(&blank_chunk).unwrap();

                    new_chunk_file
                        .write_all(encoded_chunk.as_bytes())
                        .expect("Error writing data to chunk");

                    new_chunk_file
                }
                err => panic!("Opening new file for chunk failed with: {:?}", err),
            },
        };

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

        if !self.cached_file_names.is_empty() {
            // Find the name of the previous entry
            let (previous_file_name, _) = self
                .cached_file_names
                .iter()
                .find(|(_, &array_index)| array_index == last_used_index)
                .expect("The last used index should always have a name");

            self.cached_file_names.remove(&previous_file_name.clone());
        }
        self.cached_file_names
            .insert(file_name.to_string(), last_used_index);
        // Replace the timestamp with the new timestamp
        self.last_used_times[last_used_index] = Instant::now();
        self.cached_chunk_files[last_used_index] = Some(chunk_file);

        self.cached_chunk_files[last_used_index].as_ref().unwrap()
    }

    /// `fetch_chunk_by_pos` takes in the position of a chunk, and returns the
    /// data of the chunk from disk
    ///
    /// This operation is cached, if possible, so that subsequent accesses to
    /// the same chunk are handled by the same file
    pub fn fetch_chunk_by_pos(&mut self, pos: &ChunkPos) -> Result<ChunkData, Box<dyn Error>> {
        let file_name = pos.storage_file_name();

        let file_index = self.cached_file_names.get(file_name.as_str());

        let chunk_file = match file_index {
            Some(index) => self.cached_chunk_files[*index].as_ref().unwrap(),
            None => self.load_chunk_file(pos, file_name.as_str()),
        };

        let file_contents = std::io::read_to_string(chunk_file)?;

        let read_data: ChunkData = serde_json::from_str(file_contents.as_str())?;

        // let read_data: ChunkData = serde_json::from_reader(&mut file_buffer)?;
        // let read_data = ChunkData::new(&ChunkPos { x: 0, z: 0 });

        Ok(read_data)
    }
}
