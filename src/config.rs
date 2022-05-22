use serde::{Deserialize, Serialize};
use toml::from_str;

#[derive(Deserialize, Serialize, Debug)]
pub struct Feed {
    name: String,
    url: String,
}

#[derive(Deserialize, Serialize, Debug)]
pub struct Config {
    feeds: Vec<Feed>,
}

impl Config {
    pub fn from_str(str: String) -> Result<Config, toml::de::Error> {
        from_str(&str)
    }
}
