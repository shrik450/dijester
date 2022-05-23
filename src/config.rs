use serde::{Deserialize, Serialize};
use toml::from_str;

#[derive(Deserialize, Serialize, Debug)]
pub enum FeedType {
    Atom,
    RSS,
}

#[derive(Deserialize, Serialize, Debug)]
pub struct Feed {
    pub name: String,
    pub url: String,
    pub feed_type: FeedType,
}

#[derive(Deserialize, Serialize, Debug)]
pub struct Config {
    pub feeds: Vec<Feed>,
}

impl Config {
    pub fn from_str(str: String) -> Result<Config, toml::de::Error> {
        from_str(&str)
    }
}
