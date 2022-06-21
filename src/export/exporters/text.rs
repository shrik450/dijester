use crate::core::{entry::Entry, feed::Feed};
use html2text::from_read;

use super::Exporter;

#[derive(Debug)]
pub(super) struct TextExporter();

impl TextExporter {
    fn strip_html(content: String) -> String {
        from_read(content.as_bytes(), 80)
    }

    fn generate_content_for_entry(entry: Entry) -> String {
        let mut content = format!("{}\n\n", entry.title);

        if entry.author.is_some() {
            content.push_str(&format!("By: {}\n\n", entry.author.unwrap()));
        }

        content.push_str(&TextExporter::strip_html(entry.content));

        content
    }
}

impl Exporter for TextExporter {
    fn build_single_file(&self, entry: Entry) -> Vec<u8> {
        TextExporter::generate_content_for_entry(entry)
            .as_bytes()
            .to_vec()
    }

    fn build_compiled_file(&self, entries: Vec<(Feed, Vec<Entry>)>) -> Vec<u8> {
        let mut content = String::new();

        for (feed, feed_entries) in entries {
            content.push_str(&format!("{}\n", feed.name));
            content.push_str(&"=".repeat(80));
            content.push_str("\n\n");

            for entry in feed_entries {
                content.push_str(&TextExporter::generate_content_for_entry(entry));
                content.push('\n');
                content.push_str(&"-".repeat(80));
                content.push('\n');
            }
        }

        content.as_bytes().to_vec()
    }

    fn file_extension(&self) -> String {
        "txt".to_owned()
    }
}
