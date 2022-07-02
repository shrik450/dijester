use std::collections::HashMap;

use anyhow::Result;
use chrono::{DateTime, Utc};

use super::FilterConfig;
use crate::core::{entry::Entry, feed::Feed};

async fn get_last_read_mapping(namespace: &str) -> Result<HashMap<String, DateTime<Utc>>> {
    let file_contents = tokio::fs::read(format!(".{}.dijester-history", namespace)).await?;
    let mapping = toml::from_str(&String::from_utf8(file_contents)?)?;

    Ok(mapping)
}

fn filter_read(entries: Vec<Entry>, last_read: DateTime<Utc>) -> Vec<Entry> {
    entries
        .into_iter()
        .filter(|entry| entry.published_at > last_read)
        .collect()
}

fn truncate_excess(entries: Vec<Entry>, max_entries: u8) -> Vec<Entry> {
    entries.into_iter().take(max_entries.into()).collect()
}

async fn write_history_file(feeds: Vec<String>, namespace: &str) -> Result<()> {
    let last_read: HashMap<String, DateTime<Utc>> = feeds
        .into_iter()
        .map(|feed_name| (feed_name, Utc::now()))
        .collect();
    let toml_string = toml::ser::to_string(&last_read)?;
    let file_name = format!(".{}.dijester-history", namespace);
    tokio::fs::write(file_name, toml_string).await?;
    Ok(())
}

pub(super) async fn filter_for_export(
    entries: Vec<(Feed, Vec<Entry>)>,
    filter_config: FilterConfig,
) -> Vec<(Feed, Vec<Entry>)> {
    let last_read_mapping = match get_last_read_mapping(&filter_config.namespace).await {
        Ok(mapping) => mapping,
        Err(_) => {
            log::info!(
                "Couldn't find or read dijester-history file; \
                won't filter entries based on last read."
            );
            log::info!(
                "This isn't usually a problem - the flow will create \
                a dijester-history file later."
            );
            HashMap::new()
        }
    };

    let feeds: Vec<_> = entries.iter().map(|(feed, _)| feed.name.clone()).collect();
    match write_history_file(feeds, &filter_config.namespace).await {
        Ok(()) => log::info!("Wrote history file."),
        Err(_) => {
            log::warn!("Couldn't write history file. This could be a problem for future runs!")
        }
    }

    entries
        .into_iter()
        .map(|(feed, feed_entries)| {
            log::debug!("Filtering entries for feed: {}", &feed.name);

            let final_track_read = feed.track_read.unwrap_or(filter_config.global_track_read);
            let final_max_entries = feed
                .max_new_articles
                .unwrap_or(filter_config.global_max_entries);

            let mut final_entries = feed_entries;

            if final_track_read {
                let last_read = last_read_mapping.get(&feed.name);
                final_entries = match last_read {
                    Some(time) => filter_read(final_entries, time.to_owned()),
                    None => final_entries,
                }
            }

            final_entries = truncate_excess(final_entries, final_max_entries);

            log::debug!("Left with {} entries after filtering.", final_entries.len());

            (feed, final_entries)
        })
        .filter(|(_, entries)| !entries.is_empty())
        .collect()
}
