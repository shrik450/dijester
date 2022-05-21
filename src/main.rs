use clap::Parser;
use env_logger::{Builder, Env, Target};
use log;

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
}
