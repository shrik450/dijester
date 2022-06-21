use crate::core::{
    entry::Entry,
    feed::{
        Feed,
        FeedType::{Atom, Rss},
    },
};
use anyhow::Result;
use atom_syndication::Feed as AtomFeed;
use log::trace;
use rss::Channel as RssFeed;

pub async fn fetch(feed: &Feed) -> Result<Vec<Entry>> {
    trace!("Fetching entries for feed: {}", feed.name);

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

    trace!("Fetched {} entries.", entries.len());

    Ok(entries)
}

pub async fn fetch_all(feeds: Vec<Feed>) -> Result<Vec<(Feed, Vec<Entry>)>> {
    let mut entries: Vec<(Feed, Vec<Entry>)> = vec![];

    for feed in feeds {
        let feed_entries = fetch(&feed).await?;
        entries.push((feed, feed_entries));
    }

    Ok(entries)
}
