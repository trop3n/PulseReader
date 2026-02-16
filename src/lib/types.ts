export interface Feed {
  id: number;
  type: 'rss' | 'email';
  name: string;
  url: string | null;
  config: string | null;
  last_fetched: string | null;
  created_at: string;
  total_count: number;
  unread_count: number;
}

export interface Article {
  id: number;
  feed_id: number;
  title: string;
  content: string | null;
  url: string | null;
  author: string | null;
  published_at: string | null;
  is_read: number;
  created_at: string;
  feed_name: string;
  feed_type?: string;
}
