mod config;
mod core;
mod extensions;

use clap::Parser;
use env_logger::{Builder, Env, Target};
use log::info;

#[derive(Parser)]
#[clap(author, version, about)]
struct Cli {
    /// The path to a dijester config file. For more information, read: <url-here>.
    #[clap(parse(from_os_str))]
    config_file: std::path::PathBuf,
}

#[tokio::main]
async fn main() {
    Builder::from_env(Env::default())
        .target(Target::Stdout)
        .init();

    let args = Cli::parse();

    let conf: config::Config = args.config_file.try_into().unwrap();

    info!("Successfully loaded config from file.")
}
