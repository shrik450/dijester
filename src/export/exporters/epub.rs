use crate::core::{entry::Entry, feed::Feed};

use super::Exporter;

pub(super) struct EpubExporter();

impl Exporter for EpubExporter {
    fn build_single_file(&self, entry: Entry) -> Vec<u8> {
        todo!()
    }

    fn build_compiled_file(&self, entries: Vec<(Feed, Vec<Entry>)>) -> Vec<u8> {
        todo!()
    }
}