<script lang="ts">
  const articles = [
    { id: '1', title: 'Understanding React Server Components', feed: 'Tech Blogs', time: '2h ago', unread: true },
    { id: '2', title: 'The Future of TypeScript 5.5', feed: 'Tech Blogs', time: '4h ago', unread: true },
    { id: '3', title: 'Building Modern CLI Tools', feed: 'Newsletter', time: '1d ago', unread: false },
    { id: '4', title: 'CSS Container Queries Are Here', feed: 'Tech Blogs', time: '2d ago', unread: false },
    { id: '5', title: 'Introduction to WebGPU', feed: 'Papers', time: '3d ago', unread: false },
    { id: '6', title: 'System Design Interview Tips', feed: 'Newsletter', time: '4d ago', unread: false },
  ];

  interface Props {
    onSelect?: (id: string) => void;
  }

  let { onSelect }: Props = $props();

  let selectedId = $state<string | null>(null);

  function handleSelect(id: string) {
    selectedId = id;
    onSelect?.(id);
  }
</script>

<div class="article-list">
  {#each articles as article}
    <button
      class="article-item"
      class:selected={selectedId === article.id}
      class:unread={article.unread}
      onclick={() => handleSelect(article.id)}
    >
      {#if article.unread}
        <span class="unread-dot"></span>
      {/if}
      <div class="article-content">
        <span class="article-title">{article.title}</span>
        <div class="article-meta">
          <span class="article-feed">{article.feed}</span>
          <span class="separator">·</span>
          <span class="article-time">{article.time}</span>
        </div>
      </div>
    </button>
  {:else}
    <div class="empty-state">
      <p>No articles yet</p>
    </div>
  {/each}
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

  .empty-state {
    display: flex;
    align-items: center;
    justify-content: center;
    height: 200px;
    color: var(--text-muted);
    font-size: 14px;
  }
</style>
