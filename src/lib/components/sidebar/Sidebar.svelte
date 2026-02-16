<script lang="ts">
  import { onMount } from 'svelte';
  import { feeds, totalUnreadCount, rssFeeds, feedsLoading } from '$lib/stores/feeds';
  import { articles } from '$lib/stores/articles';
  import { ui } from '$lib/stores/ui';

  let loading = $state(false);
  let activeItem = $state<string>('all');

  onMount(() => {
    feeds.load();
  });

  async function handleRefresh() {
    loading = true;
    try {
      await feeds.fetchAll();
      loadArticlesForCurrentView();
    } finally {
      loading = false;
    }
  }

  function loadArticlesForCurrentView() {
    const feedId = $ui.selectedFeedId;
    const viewMode = $ui.viewMode;
    
    if (viewMode === 'unread') {
      articles.load({ unreadOnly: true });
    } else if (feedId) {
      articles.load({ feedId });
    } else {
      articles.load({});
    }
  }

  function handleSelect(itemId: string) {
    activeItem = itemId;
    
    if (itemId === 'all') {
      ui.selectFeed(null);
      ui.setViewMode('all');
      articles.load({});
    } else if (itemId === 'unread') {
      ui.selectFeed(null);
      ui.setViewMode('unread');
      articles.load({ unreadOnly: true });
    } else if (itemId.startsWith('feed-')) {
      const feedId = parseInt(itemId.replace('feed-', ''));
      ui.selectFeed(feedId);
      ui.setViewMode('all');
      articles.load({ feedId });
    }
  }

  function formatCount(count: number): string {
    if (count >= 1000) return `${Math.floor(count / 1000)}k`;
    return String(count);
  }
</script>

<nav class="sidebar-nav">
  <div class="nav-section">
    <button
      class="nav-item"
      class:active={activeItem === 'all'}
      onclick={() => handleSelect('all')}
    >
      <span class="nav-icon">
        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <polyline points="22 12 16 12 14 15 10 15 8 12 2 12"/>
          <path d="M5.45 5.11 2 12v6a2 2 0 0 0 2 2h16a2 2 0 0 0 2-2v-6l-3.45-6.89A2 2 0 0 0 16.76 4H7.24a2 2 0 0 0-1.79 1.11z"/>
        </svg>
      </span>
      <span class="nav-label">All Items</span>
    </button>
    <button
      class="nav-item"
      class:active={activeItem === 'unread'}
      onclick={() => handleSelect('unread')}
    >
      <span class="nav-icon">
        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <circle cx="12" cy="12" r="10"/>
        </svg>
      </span>
      <span class="nav-label">Unread</span>
      {#if $totalUnreadCount > 0}
        <span class="nav-count">{formatCount($totalUnreadCount)}</span>
      {/if}
    </button>
  </div>

  <div class="nav-section">
    <div class="section-header">
      <span class="section-label">Feeds</span>
      <button class="section-add" title="Add Feed" onclick={() => ui.showAddFeedModal()}>
        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M12 5v14M5 12h14"/>
        </svg>
      </button>
    </div>
    
    {#if $feedsLoading}
      <div class="loading-state">Loading...</div>
    {:else if $rssFeeds.length === 0}
      <div class="empty-section">
        <p>No feeds yet</p>
        <button class="add-first-feed" onclick={() => ui.showAddFeedModal()}>Add your first feed</button>
      </div>
    {:else}
      {#each $rssFeeds as feed}
        <button
          class="nav-item"
          class:active={activeItem === `feed-${feed.id}`}
          onclick={() => handleSelect(`feed-${feed.id}`)}
        >
          <span class="nav-icon">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M4 11a9 9 0 0 1 9 9"/>
              <path d="M4 4a16 16 0 0 1 16 16"/>
              <circle cx="5" cy="19" r="1"/>
            </svg>
          </span>
          <span class="nav-label">{feed.name}</span>
          {#if feed.unread_count > 0}
            <span class="nav-count">{formatCount(feed.unread_count)}</span>
          {/if}
        </button>
      {/each}
    {/if}
  </div>

  <div class="sidebar-footer">
    <button class="refresh-btn" onclick={handleRefresh} disabled={loading}>
      <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class:spinning={loading}>
        <path d="M21 12a9 9 0 1 1-9-9c2.52 0 4.93 1 6.74 2.74L21 8"/>
        <path d="M21 3v5h-5"/>
      </svg>
      {loading ? 'Refreshing...' : 'Refresh All'}
    </button>
  </div>
</nav>

<style>
  .sidebar-nav {
    display: flex;
    flex-direction: column;
    height: 100%;
    padding: 8px;
  }

  .nav-section {
    margin-bottom: 16px;
  }

  .section-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 6px 8px;
    margin-bottom: 2px;
  }

  .section-label {
    font-size: 11px;
    font-weight: 600;
    text-transform: uppercase;
    letter-spacing: 0.5px;
    color: var(--text-muted);
  }

  .section-add {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 20px;
    height: 20px;
    border-radius: var(--radius-sm);
    color: var(--text-muted);
    transition: all var(--transition-fast);
  }

  .section-add:hover {
    background: var(--bg-hover);
    color: var(--text-secondary);
  }

  .nav-item {
    display: flex;
    align-items: center;
    width: 100%;
    padding: 6px 8px;
    border-radius: var(--radius-sm);
    color: var(--text-secondary);
    transition: all var(--transition-fast);
    text-align: left;
  }

  .nav-item:hover {
    background: var(--bg-hover);
    color: var(--text-primary);
  }

  .nav-item.active {
    background: var(--bg-active);
    color: var(--text-primary);
  }

  .nav-icon {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 20px;
    height: 20px;
    margin-right: 8px;
    flex-shrink: 0;
  }

  .nav-label {
    flex: 1;
    font-size: 13px;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .nav-count {
    font-size: 12px;
    color: var(--text-muted);
    background: var(--bg-tertiary);
    padding: 1px 6px;
    border-radius: 10px;
  }

  .loading-state {
    padding: 12px 8px;
    font-size: 13px;
    color: var(--text-muted);
    text-align: center;
  }

  .empty-section {
    padding: 12px 8px;
    text-align: center;
  }

  .empty-section p {
    font-size: 12px;
    color: var(--text-muted);
    margin-bottom: 8px;
  }

  .add-first-feed {
    font-size: 12px;
    color: var(--accent);
    text-decoration: underline;
  }

  .add-first-feed:hover {
    color: var(--accent-hover);
  }

  .sidebar-footer {
    margin-top: auto;
    padding-top: 8px;
    border-top: 1px solid var(--border-light);
  }

  .refresh-btn {
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 6px;
    width: 100%;
    padding: 8px 12px;
    background: var(--bg-tertiary);
    border: 1px solid var(--border);
    border-radius: var(--radius-md);
    color: var(--text-secondary);
    font-size: 13px;
    transition: all var(--transition-fast);
  }

  .refresh-btn:hover:not(:disabled) {
    background: var(--bg-hover);
    border-color: var(--accent);
    color: var(--text-primary);
  }

  .refresh-btn:disabled {
    opacity: 0.6;
    cursor: not-allowed;
  }

  .spinning {
    animation: spin 1s linear infinite;
  }

  @keyframes spin {
    from { transform: rotate(0deg); }
    to { transform: rotate(360deg); }
  }
</style>
