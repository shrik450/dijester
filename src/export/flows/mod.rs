use std::fmt;

use crate::{
    config::export_config::ExportType,
    core::{entry::Entry, feed::Feed},
};

use super::{exporters::Exporter, write_actions::WriteAction};

mod compiled;
mod flat_files;

/// An **ExportFlow** runs all actions that need to be performed to generate a
/// digest for a particular **[ExportType]**
///
/// To understand why this is necessary, consider that a flat files export could
/// require making a lot of new directories, while a compiled export will only
/// require one file. This isn't entirely captured by the [Exporter], as that
/// trait only deals with converting entries into files of a particular format.
pub(super) trait ExportFlow: fmt::Debug {
    /// Determine all writes that need to be performed to export the `entries`.
    ///
    /// Returns:
    ///
    /// A list of [WriteAction]s that need to be performed. As of now, it is
    /// expected that these are performed sequentially as they *may* be order
    /// dependent. In the future, we could instead return a DAG of writes that
    /// would allow for async writing.
    fn export(
        &self,
        name: String,
        entries: Vec<(Feed, Vec<Entry>)>,
        exporter: Box<dyn Exporter>,
    ) -> anyhow::Result<Vec<WriteAction>>;
}

pub(super) fn get_flow_for_type(export_type: ExportType) -> Box<dyn ExportFlow> {
    match export_type {
        ExportType::Compiled => Box::new(compiled::Compiled()),
        ExportType::FlatFiles => Box::new(flat_files::FlatFiles()),
    }
}
