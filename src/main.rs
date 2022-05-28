mod config;
mod entry;
mod feed;

use clap::Parser;
use config::Config;
use env_logger::{Builder, Env, Target};
use futures::future;

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

    match wrapped_main(args.config_file).await {
        Ok(_) => (),
        Err(err) => {
            log::error!("Could not read config file.");
            log::error!("{}", err.to_string());
            log::error!("Exiting.");
            std::process::exit(1)
        }
    }
}

async fn wrapped_main(config_file: std::path::PathBuf) -> anyhow::Result<()> {
    log::debug!("Loading config file: {:#?}", config_file);
    let conf = Config::try_from(config_file)?;
    let promises = conf.feeds.into_iter().map(|feed| feed.load_entries());

    future::join_all(promises).await;

    Ok(())
}
