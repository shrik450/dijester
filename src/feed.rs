use atom_syndication::Feed as AtomFeed;
use rss::Channel as RssFeed;
use serde::{Deserialize, Serialize};

#[derive(Deserialize, Serialize, Debug)]
pub enum FeedType {
    Atom,
    RSS,
}

#[derive(Deserialize, Serialize, Debug)]
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

                println!("{:#?}", first_entry)
            }
            FeedType::RSS => {
                let channel = RssFeed::read_from(&resp[..])?;
                let first_entry = channel.items.first();

                println!("{:#?}", first_entry)
            }
        }

        Ok(())
    }
}
