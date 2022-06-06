use std::fs::File;

use crate::config::ExportFormat;
use crate::entry::Entry;
use crate::export::FileOutput;
use crate::feed::Feed;

use self::html_exporter::HTMLExporter;

mod html_exporter;

pub(in crate::export) trait Exporter: std::fmt::Debug {
    fn export(&self, entry: Entry, feed: &Feed) -> FileOutput;
    /// This is probably not right.
    fn compile(&self, outputs: Vec<FileOutput>) -> FileOutput;
}

pub(in crate::export) fn get_exporter(config: ExportFormat) -> impl Exporter {
    match config {
        ExportFormat::HTML => HTMLExporter(),
        ExportFormat::MD => todo!(),
        ExportFormat::EPUB => todo!(),
    }
}
