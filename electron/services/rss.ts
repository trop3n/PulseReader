import Parser from 'rss-parser';

const parser = new Parser();

export interface FeedItem {
  title: string;
  link: string;
  content: string;
  contentSnippet: string;
  pubDate?: string;
  creator?: string;
  isoDate?: string;
}

export interface FeedData {
  title: string;
  description?: string;
  link?: string;
  items: FeedItem[];
}

export async function fetchFeed(url: string): Promise<FeedData> {
  const feed = await parser.parseURL(url);
  
  return {
    title: feed.title || 'Untitled Feed',
    description: feed.description,
    link: feed.link,
    items: feed.items.map((item) => ({
      title: item.title || 'Untitled',
      link: item.link || '',
      content: item.content || item.contentSnippet || '',
      contentSnippet: item.contentSnippet || '',
      pubDate: item.pubDate,
      creator: item.creator,
      isoDate: item.isoDate
    }))
  };
}
