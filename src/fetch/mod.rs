use crate::core::{
    entry::Entry,
    feed::{
        Feed,
        FeedType::{Atom, Rss},
    },
};
use anyhow::Result;
use atom_syndication::Feed as AtomFeed;
use rss::Channel as RssFeed;

mod filter;

pub struct FilterConfig {
    pub namespace: String,
    pub global_track_read: bool,
    pub global_max_entries: u8,
}

async fn fetch_entries_for_feed(feed: &Feed) -> Result<Vec<Entry>> {
    log::debug!("Fetching entries for feed: {}", feed.name);

    let owned_feed = feed.to_owned();
    let data = reqwest::get(owned_feed.url).await?.bytes().await?;

    let entries: Vec<Entry> = match owned_feed.feed_type {
        Atom => AtomFeed::read_from(&data[..])?
            .entries
            .into_iter()
            .map(Entry::from)
            .collect(),
        Rss => RssFeed::read_from(&data[..])?
            .items()
            .iter()
            .cloned()
            .map(Entry::from)
            .collect(),
    };

    log::debug!("Fetched {} entries.", entries.len());

    Ok(entries)
}

pub async fn fetch(
    feeds: Vec<Feed>,
    filter_config: FilterConfig,
) -> Result<Vec<(Feed, Vec<Entry>)>> {
    let mut entries: Vec<(Feed, Vec<Entry>)> = vec![];

    for feed in feeds {
        let feed_entries = fetch_entries_for_feed(&feed).await?;
        entries.push((feed, feed_entries));
    }

    let filtered_entries = filter::filter_for_export(entries, filter_config).await;

    Ok(filtered_entries)
}
