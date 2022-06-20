use crate::{
    core::{entry::Entry, feed::Feed},
    export::{
        exporters::Exporter,
        flows::ExportFlow,
        write_actions::{CreateDirectoryAction, CreateFileAction, WriteAction},
    },
};

pub(super) struct FlatFiles();

impl ExportFlow for FlatFiles {
    fn export(
        &self,
        entries: Vec<(Feed, Vec<Entry>)>,
        exporter: Box<dyn Exporter>,
    ) -> anyhow::Result<Vec<WriteAction>> {
        let a = entries
            .into_iter()
            .flat_map(|(feed, entries)| {
                let mut actions: Vec<WriteAction> = vec![];
                actions.push(WriteAction::CreateDirectory(CreateDirectoryAction {
                    relative_path: feed.name.to_owned(),
                }));
                entries
                    .into_iter()
                    .map(|entry| {
                        let relative_path = format!("{}/{}", feed.name, entry.title);
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

        Ok(a)
    }
}
