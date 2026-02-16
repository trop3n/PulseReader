<script lang="ts">
  import { onMount } from 'svelte';
  import type { Article } from '$lib/types';

  interface Props {
    article?: Article | null;
  }

  let { article }: Props = $props();

  function formatContent(content: string | null): string {
    if (!content) return '<p>No content available</p>';
    return content;
  }

  function formatDate(dateStr: string | null): string {
    if (!dateStr) return '';
    const date = new Date(dateStr);
    return date.toLocaleDateString('en-US', {
      weekday: 'long',
      year: 'numeric',
      month: 'long',
      day: 'numeric'
    });
  }

  function handleOpenOriginal() {
    if (article?.url) {
      window.open(article.url, '_blank');
    }
  }
</script>

<div class="reader">
  {#if article}
    <div class="reader-toolbar">
      <div class="toolbar-info">
        <span class="feed-name">{article.feed_name}</span>
      </div>
      <div class="toolbar-spacer"></div>
      <div class="toolbar-actions">
        {#if article.url}
          <button class="toolbar-btn" onclick={handleOpenOriginal} title="Open original">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M18 13v6a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h6"/>
              <polyline points="15 3 21 3 21 9"/>
              <line x1="10" x2="21" y1="14" y2="3"/>
            </svg>
          </button>
        {/if}
        <button class="toolbar-btn primary" title="Send to Obsidian">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/>
            <polyline points="17 8 12 3 7 8"/>
            <line x1="12" x2="12" y1="3" y2="15"/>
          </svg>
          <span>Send to Obsidian</span>
        </button>
      </div>
    </div>

    <div class="reader-content">
      <article class="article">
        <h1>{article.title}</h1>
        <div class="meta">
          {#if article.author}
            <span class="author">{article.author}</span>
            <span class="separator">·</span>
          {/if}
          <time>{formatDate(article.published_at)}</time>
        </div>

        <div class="content">
          {@html formatContent(article.content)}
        </div>
      </article>
    </div>
  {:else}
    <div class="empty-state">
      <svg width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
        <path d="M4 19.5v-15A2.5 2.5 0 0 1 6.5 2H20v20H6.5a2.5 2.5 0 0 1 0-5H20"/>
      </svg>
      <p>Select an article to read</p>
    </div>
  {/if}
</div>

<style>
  .reader {
    display: flex;
    flex-direction: column;
    flex: 1;
    height: 100%;
    overflow: hidden;
    background: var(--bg-primary);
  }

  .reader-toolbar {
    display: flex;
    align-items: center;
    padding: 8px 16px;
    background: var(--bg-secondary);
    border-bottom: 1px solid var(--border-light);
  }

  .toolbar-info {
    display: flex;
    align-items: center;
    gap: 8px;
  }

  .feed-name {
    font-size: 12px;
    font-weight: 500;
    color: var(--text-muted);
    text-transform: uppercase;
    letter-spacing: 0.5px;
  }

  .toolbar-spacer {
    flex: 1;
  }

  .toolbar-actions {
    display: flex;
    align-items: center;
    gap: 8px;
  }

  .toolbar-btn {
    display: flex;
    align-items: center;
    gap: 6px;
    padding: 6px 10px;
    border-radius: var(--radius-sm);
    color: var(--text-secondary);
    font-size: 13px;
    transition: all var(--transition-fast);
  }

  .toolbar-btn:hover {
    background: var(--bg-hover);
    color: var(--text-primary);
  }

  .toolbar-btn.primary {
    background: var(--accent);
    color: white;
  }

  .toolbar-btn.primary:hover {
    background: var(--accent-hover);
  }

  .reader-content {
    flex: 1;
    overflow-y: auto;
  }

  .article {
    max-width: 680px;
    margin: 0 auto;
    padding: 48px 24px;
  }

  .article h1 {
    font-size: 28px;
    font-weight: 600;
    color: var(--text-primary);
    margin-bottom: 12px;
    line-height: 1.3;
  }

  .article .meta {
    display: flex;
    align-items: center;
    font-size: 13px;
    color: var(--text-muted);
    margin-bottom: 32px;
  }

  .article .author {
    color: var(--text-secondary);
  }

  .article .separator {
    margin: 0 8px;
  }

  .article .content {
    font-size: 16px;
    line-height: 1.7;
    color: var(--text-primary);
  }

  .article .content :global(p) {
    margin-bottom: 1.5em;
  }

  .article .content :global(h2) {
    font-size: 22px;
    font-weight: 600;
    margin-top: 2em;
    margin-bottom: 0.75em;
  }

  .article .content :global(h3) {
    font-size: 18px;
    font-weight: 600;
    margin-top: 1.5em;
    margin-bottom: 0.5em;
  }

  .article .content :global(ul),
  .article .content :global(ol) {
    margin-bottom: 1.5em;
    padding-left: 1.5em;
  }

  .article .content :global(li) {
    margin-bottom: 0.5em;
  }

  .article .content :global(blockquote) {
    border-left: 3px solid var(--accent);
    padding-left: 1em;
    margin: 1.5em 0;
    color: var(--text-secondary);
    font-style: italic;
  }

  .article .content :global(pre) {
    background: var(--bg-secondary);
    padding: 16px;
    border-radius: var(--radius-md);
    overflow-x: auto;
    margin: 1.5em 0;
    font-size: 14px;
  }

  .article .content :global(code) {
    background: var(--bg-secondary);
    padding: 2px 6px;
    border-radius: var(--radius-sm);
    font-size: 0.9em;
  }

  .article .content :global(pre code) {
    background: none;
    padding: 0;
  }

  .article .content :global(img) {
    max-width: 100%;
    height: auto;
    border-radius: var(--radius-md);
    margin: 1.5em 0;
  }

  .article .content :global(a) {
    color: var(--text-link);
  }

  .article .content :global(a:hover) {
    text-decoration: underline;
  }

  .empty-state {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    flex: 1;
    gap: 16px;
    color: var(--text-muted);
  }

  .empty-state p {
    font-size: 14px;
  }
</style>
