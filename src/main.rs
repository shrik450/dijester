mod config;
mod feed;

use clap::Parser;
use config::Config;
use env_logger::{Builder, Env, Target};
use futures::future;

#[derive(Parser)]
#[clap(author, version, about, long_about = None)]
struct Cli {
    #[clap(parse(from_os_str))]
    config_file: std::path::PathBuf,
}

#[tokio::main]
async fn main() {
    Builder::from_env(Env::default())
        .target(Target::Stdout)
        .init();

    let args = Cli::parse();

    wrapped_main(args.config_file).await;
}

async fn wrapped_main(config_file: std::path::PathBuf) -> anyhow::Result<()> {
    log::debug!("Loading config file: {:#?}", config_file);
    let conf = Config::try_from(config_file)?;
    let promises = conf.feeds.into_iter().map(|feed| feed.load_entries());

    future::join_all(promises).await;

    Ok(())
}

fn report_read_error_and_exit(err: Box<dyn std::error::Error>) {
    log::error!("Could not read config file: {}", err.to_string());
    log::error!("Exiting.");
    std::process::exit(1)
}
