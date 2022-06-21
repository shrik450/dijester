use crate::{
    core::{entry::Entry, feed::Feed},
    export::{
        exporters::Exporter,
        write_actions::{CreateFileAction, WriteAction},
    },
};

use super::ExportFlow;

#[derive(Debug)]
pub(super) struct Compiled();

impl ExportFlow for Compiled {
    fn export(
        &self,
        name: String,
        entries: Vec<(Feed, Vec<Entry>)>,
        exporter: Box<dyn Exporter>,
    ) -> anyhow::Result<Vec<WriteAction>> {
        let content = exporter.build_compiled_file(entries);

        Ok(vec![WriteAction::CreateFile(CreateFileAction {
            relative_path: format!("{}.{}", name, exporter.file_extension()),
            content,
        })])
    }
}
