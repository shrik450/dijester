use log::{debug, info};

use crate::{
    config::export_config::ExportConfig,
    core::{entry::Entry, feed::Feed},
};

mod exporters;
mod flows;
mod write_actions;

pub async fn export(
    name: String,
    entries: Vec<(Feed, Vec<Entry>)>,
    export_config: ExportConfig,
) -> anyhow::Result<()> {
    let exporter = exporters::get_exporter_for_format(export_config.export_format);
    let flow = flows::get_flow_for_type(export_config.export_type);
    let current_time = chrono::Local::now().to_rfc2822();

    let export_name = format!("{} - {}", name, current_time);

    debug!("Loaded exporter {:#?} for flow {:#?}", exporter, flow);
    info!("Exporting to: {}", export_name);

    let actions = flow.export(export_name, entries, exporter)?;

    // We can't just join them all and await the joined promise - we're
    // promising ordered execution of writes right now.
    for action in actions {
        action.write(&export_config.destination).await?
    }

    info!("Export completed.");

    Ok(())
}
