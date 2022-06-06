use std::path::PathBuf;

use crate::config::ExportConfig;
use crate::entry::Entry;
use crate::extensions::JoinableIterator;
use crate::feed::Feed;

mod exporter;

use anyhow::Ok;
use exporter::Exporter;

pub(in crate::export) struct FolderOutput {
    relative_path: String,
}

impl FolderOutput {
    async fn write(self, base_path: &PathBuf) -> std::io::Result<()> {
        tokio::fs::create_dir(base_path.join(self.relative_path)).await
    }
}

pub(in crate::export) struct FileOutput {
    relative_path: String,
    content: String,
}

impl FileOutput {
    async fn write(self, base_path: &PathBuf) -> std::io::Result<()> {
        tokio::fs::write(base_path.join(self.relative_path), self.content).await
    }
}

pub async fn export(entries: Vec<(Feed, Vec<Entry>)>, config: ExportConfig) -> anyhow::Result<()> {
    let exp = exporter::get_exporter(config.export_format);
    log::debug!("Using exporter {:#?}.", exp);

    let mut folder_outputs: Vec<FolderOutput> = Vec::new();
    let mut file_outputs: Vec<FileOutput> = Vec::new();

    for (feed, feed_entries) in entries {
        // Recipe for a shotgun change alright
        folder_outputs.push(FolderOutput {
            relative_path: feed.name.clone(),
        });

        let mut feed_file_outputs: Vec<_> = feed_entries
            .into_iter()
            .map(|entry| exp.export(entry, &feed))
            .collect();

        file_outputs.append(&mut feed_file_outputs);
    }

    // We've got to create all the parent folders first!
    folder_outputs
        .into_iter()
        .map(|o| o.write(&config.destination))
        .join_all()
        .await
        .into_iter()
        .for_each(|res| log::trace!("{:?}", res));

    file_outputs
        .into_iter()
        .map(|o| o.write(&config.destination))
        .join_all()
        .await
        .into_iter()
        .for_each(|res| log::trace!("{:?}", res));

    log::debug!("Wrote all files.");

    Ok(())
}
