use crate::{
    core::{entry::Entry, feed::Feed},
    export::{
        exporters::Exporter,
        flows::ExportFlow,
        write_actions::{CreateDirectoryAction, CreateFileAction, WriteAction},
    },
};

#[derive(Debug)]
pub(super) struct FlatFiles();

impl ExportFlow for FlatFiles {
    fn export(
        &self,
        name: String,
        entries: Vec<(Feed, Vec<Entry>)>,
        exporter: Box<dyn Exporter>,
    ) -> anyhow::Result<Vec<WriteAction>> {
        let all_actions = entries
            .into_iter()
            .flat_map(|(feed, entries)| {
                let mut actions: Vec<WriteAction> = vec![];
                let relative_path = format!("{}/{}", name, feed.name);
                actions.push(WriteAction::CreateDirectory(CreateDirectoryAction {
                    relative_path,
                }));

                entries
                    .into_iter()
                    .map(|entry| {
                        let relative_path = format!(
                            "{}/{}/{}.{}",
                            name,
                            feed.name,
                            entry.title,
                            exporter.file_extension()
                        );
                        let content = exporter.build_single_file(entry);

                        CreateFileAction {
                            relative_path,
                            content,
                        }
                    })
                    .for_each(|action| actions.push(WriteAction::CreateFile(action)));

                actions
            })
            .collect();

        Ok(all_actions)
    }
}
