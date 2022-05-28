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
    pub async fn load_entries(self) -> anyhow::Result<()> {
        let resp = reqwest::get(self.url).await?.bytes().await?;

        match self.feed_type {
            FeedType::Atom => {
                let channel = AtomFeed::read_from(&resp[..])?;
                let first_entry = channel.entries.first();
                let entry = Entry::from(first_entry.unwrap().to_owned());

                println!("{:#?}", entry)
            }

            FeedType::RSS => {
                let channel = RssFeed::read_from(&resp[..])?;
                let first_entry = channel.items.first();
                let entry = Entry::from(first_entry.unwrap().to_owned());

                println!("{:#?}", entry)
            }
        }

        Ok(())
    }
}
