use serde::{Deserialize, Serialize};

#[derive(Deserialize, Serialize, Debug, Clone)]
pub enum FeedType {
    Atom,
    Rss,
}

#[derive(Deserialize, Serialize, Debug, Clone)]
pub struct Feed {
    pub name: String,
    pub url: String,
    pub feed_type: FeedType,
    pub use_readability: Option<bool>,
    pub track_read: Option<bool>,
    pub max_new_articles: Option<u8>,
}
