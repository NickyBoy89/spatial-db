mod storage;
mod storage_server;
use clap::Parser;

#[derive(Parser, Debug)]
#[command(author, version, about, long_about = None)]
struct Args {
    /// `http` starts a build-in http server that listens for reads and writes
    /// to the database
    #[arg(long)]
    http: bool,
}

fn main() {
    let args = Args::parse();

    if args.http {
        println!("Proxy was enabled");
        storage_server::main();
    }
    println!("Hello, world!");
}
