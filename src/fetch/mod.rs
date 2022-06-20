use crate::core::{
    entry::Entry,
    feed::{
        Feed,
        FeedType::{Atom, RSS},
    },
};
use anyhow::Result;
use atom_syndication::Feed as AtomFeed;
use rss::Channel as RssFeed;

pub async fn fetch(feed: &Feed) -> Result<Vec<Entry>> {
    let owned_feed = feed.to_owned();
    let data = reqwest::get(owned_feed.url).await?.bytes().await?;

    let entries = match owned_feed.feed_type {
        Atom => AtomFeed::read_from(&data[..])?
            .entries
            .into_iter()
            .map(|e| Entry::from(e))
            .collect(),
        RSS => RssFeed::read_from(&data[..])?
            .items()
            .to_owned()
            .into_iter()
            .map(|e| Entry::from(e))
            .collect(),
    };

    Ok(entries)
}
