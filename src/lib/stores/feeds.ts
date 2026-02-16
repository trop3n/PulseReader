import { writable, derived, get, readonly } from 'svelte/store';
import type { Feed } from '../types';

declare global {
  interface Window {
    electronAPI: {
      feeds: {
        list: () => Promise<Feed[]>;
        add: (name: string, url: string) => Promise<Feed>;
        remove: (feedId: number) => Promise<boolean>;
        fetch: (feedId: number) => Promise<{ success: boolean; count: number }>;
        fetchAll: () => Promise<{ feedId: number; success: boolean; count?: number; error?: string }[]>;
      };
      articles: {
        list: (options?: { feedId?: number; unreadOnly?: boolean; limit?: number; offset?: number }) => Promise<any[]>;
        get: (articleId: number) => Promise<any>;
        markRead: (articleId: number) => Promise<boolean>;
        markAllRead: (feedId?: number) => Promise<boolean>;
        unreadCount: (feedId?: number) => Promise<number>;
      };
    };
  }
}

function createFeedsStore() {
  const feeds = writable<Feed[]>([]);
  const isLoading = writable(false);
  const errorStore = writable<string | null>(null);

  async function load() {
    isLoading.set(true);
    errorStore.set(null);
    try {
      const result = await window.electronAPI.feeds.list();
      feeds.set(result);
    } catch (e) {
      errorStore.set((e as Error).message);
    } finally {
      isLoading.set(false);
    }
  }

  async function add(name: string, url: string) {
    const feed = await window.electronAPI.feeds.add(name, url);
    feeds.update(f => [...f, { ...feed, total_count: 0, unread_count: 0 }]);
    return feed;
  }

  async function remove(feedId: number) {
    await window.electronAPI.feeds.remove(feedId);
    feeds.update(f => f.filter(feed => feed.id !== feedId));
  }

  async function fetch(feedId: number) {
    const result = await window.electronAPI.feeds.fetch(feedId);
    await load();
    return result;
  }

  async function fetchAll() {
    const results = await window.electronAPI.feeds.fetchAll();
    await load();
    return results;
  }

  return {
    subscribe: feeds.subscribe,
    isLoading,
    error: readonly(errorStore),
    load,
    add,
    remove,
    fetch,
    fetchAll
  };
}

export const feeds = createFeedsStore();
export const feedsLoading = derived(feeds.isLoading, $loading => $loading);

export const totalUnreadCount = derived(feeds, $feeds => 
  $feeds.reduce((sum, feed) => sum + feed.unread_count, 0)
);

export const rssFeeds = derived(feeds, $feeds => 
  $feeds.filter(f => f.type === 'rss')
);

export const emailFeeds = derived(feeds, $feeds => 
  $feeds.filter(f => f.type === 'email')
);
