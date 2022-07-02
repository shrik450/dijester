mod config;
mod core;
mod export;
mod extensions;
mod fetch;

use clap::Parser;
use env_logger::{Builder, Env, Target};
use fetch::FilterConfig;

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
    log::info!("Successfully loaded config from file.");

    let filter_config = FilterConfig {
        namespace: conf.name.clone(),
        global_track_read: conf.track_read,
        global_max_entries: conf.max_new_entries,
    };
    let entries = fetch::fetch(conf.feeds, filter_config).await?;
    log::info!("Fetched entries to export.");

    if entries.is_empty() {
        log::info!("No entries to export; won't create a digest for this run.");
    } else {
        export::export(conf.name, entries, conf.export_options).await?;
    }
    log::info!("Exported entries, exiting.");

    Ok(())
}
