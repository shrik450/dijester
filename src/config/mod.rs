use crate::core::feed::Feed;
use serde::{Deserialize, Serialize};

use self::{
    config_parse_error::ConfigParseError, export_config::ExportConfig,
    readability_config::ReadabilityConfig,
};

pub mod config_parse_error;
pub mod export_config;
pub mod readability_config;

/// A variant of [Config] that exists to allow for more generous parsing of the
/// TOML input.
///
/// We allow several config options to be ignored and default to safe values.
/// But if we use Options, that means we have to unwrap those every where. This
/// struct is what the TOML deserializes to - and then we create a safe [Config]
/// struct with the defaults substituted for missing values.
#[derive(Deserialize, Serialize, Debug)]
struct SerializableConfig {
    /// The name of this feed. This is used to name the digest.
    pub name: String,
    /// Whether you want to track which entries have been read or just pick
    /// the latest n every time. Defaults to false.
    pub track_read: Option<bool>,
    /// Global setting for at most how many new entries of a feed should be
    /// included in the digest. Defaults to 25.
    pub max_new_entries: Option<u8>,
    /// See [ExportConfig].
    pub export_options: ExportConfig,
    /// See [ReadabilityConfig]. Not required if you aren't using readability.
    pub readability_config: Option<ReadabilityConfig>,
    /// An array of feeds you want to generate digests from. See [Feed].
    pub feeds: Vec<Feed>,
}

#[derive(Debug)]
pub struct Config {
    pub name: String,
    pub track_read: bool,
    pub max_new_entries: u8,
    pub export_options: ExportConfig,
    pub readability_config: ReadabilityConfig,
    pub feeds: Vec<Feed>,
}

impl From<SerializableConfig> for Config {
    fn from(serializable_config: SerializableConfig) -> Self {
        Config {
            name: serializable_config.name,
            track_read: serializable_config.track_read.unwrap_or(false),
            max_new_entries: serializable_config.max_new_entries.unwrap_or(25),
            export_options: serializable_config.export_options,
            readability_config: serializable_config.readability_config.unwrap_or_default(),
            feeds: serializable_config.feeds,
        }
    }
}

impl TryFrom<String> for Config {
    type Error = ConfigParseError;

    fn try_from(value: String) -> Result<Self, Self::Error> {
        match toml::from_str::<SerializableConfig>(&value) {
            Ok(conf) => {
                log::trace!(
                    "Successfully deserialized content as SerializableConfig: {:#?}",
                    conf
                );
                let final_conf: Config = conf.into();
                log::trace!(
                    "Final config after substituting defaults: {:#?}",
                    final_conf
                );
                Ok(final_conf)
            }
            Err(err) => Err(ConfigParseError::ParseError(err.to_string())),
        }
    }
}

impl TryFrom<std::path::PathBuf> for Config {
    type Error = ConfigParseError;

    fn try_from(value: std::path::PathBuf) -> Result<Self, Self::Error> {
        match std::fs::read_to_string(value) {
            Ok(content) => {
                log::trace!("Successfully read file.");
                Self::try_from(content)
            }
            Err(err) => Err(ConfigParseError::FileError(err.to_string())),
        }
    }
}
