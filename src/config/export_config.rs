use serde::{Deserialize, Serialize};

/// Determines whether the output of the digest is one `Compiled` file, or
/// several `FlatFiles` files. The second option is useful if you would like to
/// further process to the output of dijester.
#[derive(Deserialize, Serialize, Debug)]
pub enum ExportType {
    Compiled,
    FlatFiles,
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
    /// The root folder to export digests to.
    pub destination: std::path::PathBuf,
    /// See [export_type::ExportType].
    pub export_type: ExportType,
    /// See [ExportFormat].
    pub export_format: ExportFormat,
}
