mod config;
mod core;
mod export;
mod extensions;
mod fetch;

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
async fn main() -> anyhow::Result<()> {
    Builder::from_env(Env::default())
        .target(Target::Stdout)
        .init();

    let args = Cli::parse();

    let conf: config::Config = args.config_file.try_into().unwrap();

    info!("Successfully loaded config from file.");

    let entries = fetch::fetch(&conf.feeds[0]).await?;

    export::export(vec![(conf.feeds[0].clone(), entries)], conf.export_options).await?;

    Ok(())
}
