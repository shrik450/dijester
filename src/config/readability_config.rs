use serde::{Deserialize, Serialize};

/// Configuration for accessing a Readability Service.
#[derive(Deserialize, Serialize, Default, Debug)]
pub struct ReadabilityConfig {
    /// The address of the Readability service.
    pub host: String,
    /// The port.
    pub port: i32,
    /// A sub path. Set / if not required.
    pub path: String,
}
