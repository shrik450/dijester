use crate::{
    entry::Entry,
    export::{exporter::Exporter, FileOutput},
    feed::Feed,
};

#[derive(Debug)]
pub struct HTMLExporter();

impl Exporter for HTMLExporter {
    fn export(&self, entry: Entry, feed: &Feed) -> FileOutput {
        FileOutput {
            relative_path: format!("{}/{}", feed.name, entry.title),
            content: entry.content.unwrap_or_else(|| "".to_string()),
        }
    }

    fn compile(&self, outputs: Vec<FileOutput>) -> FileOutput {
        todo!()
    }
}
