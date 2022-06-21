use std::fmt::Debug;

use crate::{
    config::export_config::ExportFormat,
    core::{entry::Entry, feed::Feed},
};

mod epub;
mod markdown;
mod text;

/// An **Exporter** builds files from one or many [Entr(ies)](Entry) for the
/// desired **[ExportFormat]**.
pub(super) trait Exporter: Debug {
    /// Builds an entry into a single file.
    fn build_single_file(&self, entry: Entry) -> Vec<u8>;

    /// Builds a mapping of feeds to their entries into a compiled file.
    fn build_compiled_file(&self, entries: Vec<(Feed, Vec<Entry>)>) -> Vec<u8>;

    fn file_extension(&self) -> String;
}

pub(super) fn get_exporter_for_format(format: ExportFormat) -> Box<dyn Exporter> {
    match format {
        ExportFormat::Txt => Box::new(text::TextExporter()),
        ExportFormat::Md => Box::new(markdown::MarkdownExporter()),
        ExportFormat::Epub => Box::new(epub::EpubExporter()),
    }
}
