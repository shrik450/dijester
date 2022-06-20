use crate::{
    config::export_config::ExportConfig,
    core::{entry::Entry, feed::Feed},
};

mod exporters;
mod flows;
mod write_actions;

pub async fn export(
    entries: Vec<(Feed, Vec<Entry>)>,
    export_config: ExportConfig,
) -> anyhow::Result<()> {
    let exporter = exporters::get_exporter_for_format(export_config.export_format);
    let flow = flows::get_flow_for_type(export_config.export_type);

    let actions = flow.export(entries, exporter)?;

    // We can't just join them all and await the joined promise - we're
    // promising ordered execution of the writes right now.
    for action in actions {
        action.write(&export_config.destination).await?
    }

    Ok(())
}
