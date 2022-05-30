use crate::entry::Entry;
use atom_syndication::Feed as AtomFeed;
use rss::Channel as RssFeed;
use serde::{Deserialize, Serialize};

#[derive(Deserialize, Serialize, Debug, Clone)]
pub enum FeedType {
    Atom,
    RSS,
}

#[derive(Deserialize, Serialize, Debug, Clone)]
pub struct Feed {
    pub name: String,
    pub url: String,
    pub feed_type: FeedType,
}

impl Feed {
    pub async fn load_entries(&self) -> anyhow::Result<Vec<Entry>> {
        let resp = reqwest::get(&self.url).await?.bytes().await?;

        let entries: Vec<Entry> = match self.feed_type {
            FeedType::Atom => {
                let channel = AtomFeed::read_from(&resp[..])?;
                channel
                    .entries
                    .into_iter()
                    .map(|feed_entry| Entry::from(feed_entry))
                    .collect()
            }

            FeedType::RSS => {
                let channel = RssFeed::read_from(&resp[..])?;
                channel
                    .items
                    .into_iter()
                    .map(|feed_entry| Entry::from(feed_entry))
                    .collect()
            }
        };

        Ok(entries)
    }
}
