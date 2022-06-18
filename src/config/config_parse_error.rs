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
