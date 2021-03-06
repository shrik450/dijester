use atom_syndication::Entry as AtomEntry;
use chrono::{DateTime, Utc};
use rss::Item as RssEntry;

#[derive(Debug)]
pub struct Entry {
    pub id: String,
    pub title: String,
    pub author: Option<String>,
    pub link: Option<String>,
    pub content: String,
    pub published_at: DateTime<Utc>,
}

impl From<AtomEntry> for Entry {
    fn from(value: AtomEntry) -> Entry {
        Entry {
            id: value.id,
            title: value.title.to_string(),
            author: value.authors.first().map(|author| author.name.to_owned()),
            link: value.links.first().map(|link| link.href.to_owned()),
            content: value
                .content
                .and_then(|content| content.value)
                .unwrap_or_default(),
            published_at: value.published.unwrap().to_owned().with_timezone(&Utc),
        }
    }
}

impl From<RssEntry> for Entry {
    fn from(value: RssEntry) -> Entry {
        Entry {
            id: value.guid.unwrap().value,
            title: value
                .title
                .unwrap_or_else(|| "Untitled Article".to_string()),
            author: value.author,
            link: value.link,
            content: value.content.unwrap_or_default(),
            published_at: DateTime::parse_from_rfc2822(value.pub_date.as_ref().unwrap())
                .unwrap()
                .with_timezone(&Utc),
        }
    }
}
