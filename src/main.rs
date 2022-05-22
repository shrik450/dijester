mod config;

use clap::Parser;
use config::Config;
use env_logger::{Builder, Env, Target};
use std::fs::read_to_string;

#[derive(Parser)]
#[clap(author, version, about, long_about = None)]
struct Cli {
    #[clap(parse(from_os_str))]
    config_file: std::path::PathBuf,
}

fn main() {
    Builder::from_env(Env::default())
        .target(Target::Stdout)
        .init();

    let args = Cli::parse();

    log::info!("Loading config file: {:#?}", args.config_file);
    let res = read_to_string(args.config_file);

    match res {
        Ok(content) => {
            log::info!("Opened config file.");
            parse_config_file(content)
        }
        Err(err) => report_read_error_and_exit(&err),
    }
}

fn parse_config_file(content: String) {
    let res = Config::from_str(content);

    match res {
        Ok(conf) => println!("{:#?}", conf),
        Err(err) => report_read_error_and_exit(&err),
    }
}

fn report_read_error_and_exit(err: &(dyn std::error::Error)) {
    log::error!("Could not read config file: {}", err.to_string());
    log::error!("Exiting.");
    std::process::exit(1)
}
