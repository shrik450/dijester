use crate::{
    core::{entry::Entry, feed::Feed},
    export::{exporters::Exporter, flows::ExportFlow, write_actions::WriteAction},
};

pub(super) struct FlatFiles();

impl ExportFlow for FlatFiles {
    fn export(
        &self,
        entries: Vec<(Feed, Vec<Entry>)>,
        exporter: Box<dyn Exporter>,
    ) -> anyhow::Result<Vec<WriteAction>> {
        Ok(vec![])
    }
}
