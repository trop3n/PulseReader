import { writable, derived, get } from 'svelte/store';
import type { Article } from '../types';

interface ArticlesState {
  items: Article[];
  loading: boolean;
  error: string | null;
  hasMore: boolean;
  totalUnread: number;
}

function createArticlesStore() {
  const initialState: ArticlesState = {
    items: [],
    loading: false,
    error: null,
    hasMore: true,
    totalUnread: 0
  };

  const { subscribe, set, update } = writable<ArticlesState>(initialState);

  async function load(options: { feedId?: number; unreadOnly?: boolean; reset?: boolean } = {}) {
    const { feedId, unreadOnly, reset = true } = options;
    
    update(state => ({ ...state, loading: true, error: null }));
    
    try {
      const currentItems = reset ? [] : get({ subscribe }).items;
      const offset = currentItems.length;
      
      const articles = await window.electronAPI.articles.list({
        feedId,
        unreadOnly,
        limit: 50,
        offset: reset ? 0 : offset
      });

      const totalUnread = await window.electronAPI.articles.unreadCount(feedId);

      update(state => ({
        ...state,
        items: reset ? articles : [...currentItems, ...articles],
        loading: false,
        hasMore: articles.length === 50,
        totalUnread
      }));
    } catch (e) {
      update(state => ({ ...state, loading: false, error: (e as Error).message }));
    }
  }

  async function getArticle(articleId: number): Promise<Article | null> {
    try {
      return await window.electronAPI.articles.get(articleId);
    } catch {
      return null;
    }
  }

  async function markRead(articleId: number) {
    await window.electronAPI.articles.markRead(articleId);
    update(state => ({
      ...state,
      items: state.items.map(a => 
        a.id === articleId ? { ...a, is_read: 1 } : a
      ),
      totalUnread: Math.max(0, state.totalUnread - 1)
    }));
  }

  async function markAllRead(feedId?: number) {
    await window.electronAPI.articles.markAllRead(feedId);
    update(state => ({
      ...state,
      items: state.items.map(a => ({ ...a, is_read: 1 })),
      totalUnread: 0
    }));
  }

  function reset() {
    set(initialState);
  }

  return {
    subscribe,
    load,
    getArticle,
    markRead,
    markAllRead,
    reset
  };
}

export const articles = createArticlesStore();
