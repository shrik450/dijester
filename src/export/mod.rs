use crate::config::ExportConfig;
use crate::entry::Entry;

pub async fn export(
    feed_name: String,
    entries: Vec<Entry>,
    config: &ExportConfig,
) -> anyhow::Result<()> {
    println!(
        "Here is where I'd export {} with {} entries using {:#?}",
        feed_name,
        entries.len(),
        config
    );

    Ok(())
}
