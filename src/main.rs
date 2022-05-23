mod config;

use atom_syndication::Feed as AtomFeed;
use clap::Parser;
use config::{Config, Feed, FeedType};
use env_logger::{Builder, Env, Target};
use rss::Channel;
use std::fs::read_to_string;

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

    log::debug!("Loading config file: {:#?}", args.config_file);
    let res = read_to_string(args.config_file);

    match res {
        Ok(content) => {
            log::debug!("Opened config file.");
            parse_config_file(content).await
        }
        Err(err) => report_read_error_and_exit(&err),
    }
}

async fn parse_config_file(content: String) {
    let res = Config::from_str(content);

    match res {
        Ok(conf) => {
            log::debug!("Reading feeds...");
            read_feeds(conf.feeds).await;
            ()
        }
        Err(err) => report_read_error_and_exit(&err),
    }
}

async fn read_feeds(feeds: Vec<Feed>) -> Result<(), Box<dyn std::error::Error>> {
    for feed in feeds {
        log::debug!("Reading feed: #{:#?}", feed);
        let resp = reqwest::get(feed.url).await?.bytes().await?;
        match feed.feed_type {
            FeedType::Atom => {
                let channel = AtomFeed::read_from(&resp[..]);
                println!("{:#?}", channel)
            }
            FeedType::RSS => {
                let channel = Channel::read_from(&resp[..]);
                println!("{:#?}", channel)
            }
        }
    }

    Ok(())
}

fn report_read_error_and_exit(err: &(dyn std::error::Error)) {
    log::error!("Could not read config file: {}", err.to_string());
    log::error!("Exiting.");
    std::process::exit(1)
}
