use crate::core::{entry::Entry, feed::Feed};

use super::Exporter;

pub(super) struct TextExporter();

impl Exporter for TextExporter {
    fn build_single_file(&self, entry: Entry) -> Vec<u8> {
        entry.content.unwrap().as_bytes().to_vec()
    }

    fn build_compiled_file(&self, entries: Vec<(Feed, Vec<Entry>)>) -> Vec<u8> {
        todo!()
    }
}
