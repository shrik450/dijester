use crate::{
    core::{entry::Entry, feed::Feed},
    export::{exporters::Exporter, write_actions::WriteAction},
};

use super::ExportFlow;

pub(super) struct Compiled();

impl ExportFlow for Compiled {
    fn export(
        &self,
        entries: Vec<(Feed, Vec<Entry>)>,
        exporter: Box<dyn Exporter>,
    ) -> anyhow::Result<Vec<WriteAction>> {
        todo!()
    }
}
