<script lang="ts">
  import { articles } from '$lib/stores/articles';
  import type { Article } from '$lib/types';
  import { onMount } from 'svelte';

  interface Props {
    onSelect?: (article: Article) => void;
  }

  let { onSelect }: Props = $props();

  let selectedId = $state<number | null>(null);
  
  let items = $derived($articles.items);
  let loading = $derived($articles.loading);
  let error = $derived($articles.error);

  onMount(() => {
    articles.load({});
  });

  function handleSelect(article: Article) {
    selectedId = article.id;
    onSelect?.(article);
    
    if (!article.is_read) {
      articles.markRead(article.id);
    }
  }

  function formatTime(dateStr: string | null): string {
    if (!dateStr) return '';
    
    const date = new Date(dateStr);
    const now = new Date();
    const diffMs = now.getTime() - date.getTime();
    const diffMins = Math.floor(diffMs / 60000);
    const diffHours = Math.floor(diffMs / 3600000);
    const diffDays = Math.floor(diffMs / 86400000);

    if (diffMins < 60) return `${diffMins}m ago`;
    if (diffHours < 24) return `${diffHours}h ago`;
    if (diffDays < 7) return `${diffDays}d ago`;
    
    return date.toLocaleDateString('en-US', { month: 'short', day: 'numeric' });
  }
</script>

<div class="article-list">
  {#if loading && items.length === 0}
    <div class="loading-state">
      <div class="spinner"></div>
      <p>Loading articles...</p>
    </div>
  {:else if error}
    <div class="error-state">
      <p>{error}</p>
      <button onclick={() => articles.load({})}>Retry</button>
    </div>
  {:else if items.length === 0}
    <div class="empty-state">
      <svg width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
        <path d="M4 19.5v-15A2.5 2.5 0 0 1 6.5 2H20v20H6.5a2.5 2.5 0 0 1 0-5H20"/>
      </svg>
      <p>No articles yet</p>
      <span>Add some feeds to get started</span>
    </div>
  {:else}
    {#each items as article (article.id)}
      <button
        class="article-item"
        class:selected={selectedId === article.id}
        class:unread={!article.is_read}
        onclick={() => handleSelect(article)}
      >
        {#if !article.is_read}
          <span class="unread-dot"></span>
        {/if}
        <div class="article-content">
          <span class="article-title">{article.title}</span>
          <div class="article-meta">
            <span class="article-feed">{article.feed_name}</span>
            <span class="separator">·</span>
            <span class="article-time">{formatTime(article.published_at)}</span>
          </div>
        </div>
      </button>
    {/each}

    {#if loading}
      <div class="loading-more">
        <div class="spinner small"></div>
      </div>
    {/if}
  {/if}
</div>

<style>
  .article-list {
    flex: 1;
    overflow-y: auto;
  }

  .article-item {
    display: flex;
    align-items: flex-start;
    width: 100%;
    padding: 12px 16px;
    text-align: left;
    border-bottom: 1px solid var(--border-light);
    transition: background var(--transition-fast);
  }

  .article-item:hover {
    background: var(--bg-secondary);
  }

  .article-item.selected {
    background: var(--bg-tertiary);
  }

  .unread-dot {
    width: 8px;
    height: 8px;
    margin-top: 6px;
    margin-right: 10px;
    border-radius: 50%;
    background: var(--accent);
    flex-shrink: 0;
  }

  .article-item:not(.unread) .unread-dot {
    display: none;
  }

  .article-content {
    flex: 1;
    min-width: 0;
  }

  .article-title {
    display: block;
    font-size: 13px;
    font-weight: 500;
    color: var(--text-primary);
    line-height: 1.4;
    margin-bottom: 4px;
  }

  .article-item:not(.unread) .article-title {
    color: var(--text-secondary);
  }

  .article-meta {
    display: flex;
    align-items: center;
    font-size: 12px;
    color: var(--text-muted);
  }

  .separator {
    margin: 0 4px;
  }

  .loading-state,
  .error-state,
  .empty-state {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    height: 200px;
    padding: 24px;
    text-align: center;
  }

  .loading-state p,
  .error-state p,
  .empty-state p {
    margin-top: 12px;
    color: var(--text-secondary);
    font-size: 14px;
  }

  .empty-state span {
    font-size: 12px;
    color: var(--text-muted);
  }

  .error-state button {
    margin-top: 12px;
    padding: 6px 12px;
    background: var(--bg-tertiary);
    border: 1px solid var(--border);
    border-radius: var(--radius-sm);
    color: var(--text-secondary);
    font-size: 12px;
  }

  .error-state button:hover {
    background: var(--bg-hover);
  }

  .spinner {
    width: 24px;
    height: 24px;
    border: 2px solid var(--border);
    border-top-color: var(--accent);
    border-radius: 50%;
    animation: spin 0.8s linear infinite;
  }

  .spinner.small {
    width: 16px;
    height: 16px;
    border-width: 2px;
  }

  .loading-more {
    display: flex;
    justify-content: center;
    padding: 16px;
  }

  @keyframes spin {
    to { transform: rotate(360deg); }
  }
</style>
