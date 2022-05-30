use crate::feed::Feed;
use serde::{Deserialize, Serialize};

#[derive(Debug)]
pub enum ConfigParseError {
    FileError(String),
    ParseError(String),
}

impl std::fmt::Display for ConfigParseError {
    fn fmt(&self, fmt: &mut std::fmt::Formatter) -> std::fmt::Result {
        match self {
            ConfigParseError::FileError(str) => {
                let error_string =
                    format!("Failed to parse: Could not load file.\n Caused by: {}", str);
                fmt.write_str(&error_string)
            }
            ConfigParseError::ParseError(str) => {
                let error_string = format!(
                    "Failed to parse: Could not parse TOML as Feed.\n Caused by: {}",
                    str
                );
                fmt.write_str(&error_string)
            }
        }
    }
}

impl std::error::Error for ConfigParseError {}

/// Determines whether the output of the digest is one `Compiled` file, or
/// several `Separate` files. The second option is useful if you would like to
/// further process to the output of dijester.
#[derive(Deserialize, Serialize, Debug)]
pub enum ExportType {
    Compiled,
    Separate,
}

/// Determines the file format of the output.
#[derive(Deserialize, Serialize, Debug)]
pub enum ExportFormat {
    HTML,
    MD,
    EPUB,
}

#[derive(Deserialize, Serialize, Debug)]
pub struct ExportConfig {
    pub destination: std::path::PathBuf,
    pub export_type: ExportType,
    pub export_format: ExportFormat,
}

#[derive(Deserialize, Serialize, Debug)]
pub struct Config {
    /// The name of this feed. This is used to name the digest.
    pub name: String,
    /// How you'd like the digest to be stored.
    pub export_options: ExportConfig,
    /// An array of feeds you want to generate digests from.
    pub feeds: Vec<Feed>,
}

impl TryFrom<String> for Config {
    type Error = ConfigParseError;

    fn try_from(value: String) -> Result<Self, Self::Error> {
        match toml::from_str::<Config>(&value) {
            Ok(conf) => Ok(conf),
            Err(err) => Err(ConfigParseError::ParseError(err.to_string())),
        }
    }
}

impl TryFrom<std::path::PathBuf> for Config {
    type Error = ConfigParseError;

    fn try_from(value: std::path::PathBuf) -> Result<Self, Self::Error> {
        match std::fs::read_to_string(value) {
            Ok(content) => Self::try_from(content),
            Err(err) => Err(ConfigParseError::FileError(err.to_string())),
        }
    }
}
