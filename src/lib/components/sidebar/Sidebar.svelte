<script lang="ts">
  const navItems = [
    { id: 'all', label: 'All Items', icon: 'inbox', count: 0 },
    { id: 'unread', label: 'Unread', icon: 'circle', count: 0 },
  ];
  
  const feedSections = [
    {
      label: 'Feeds',
      items: [
        { id: 'rss-tech', label: 'Tech Blogs', icon: 'rss', count: 12 },
        { id: 'rss-news', label: 'News', icon: 'rss', count: 5 },
      ]
    },
    {
      label: 'Newsletters',
      items: [
        { id: 'email-tech', label: 'Tech Newsletters', icon: 'mail', count: 3 },
      ]
    },
    {
      label: 'PDFs',
      items: [
        { id: 'pdf-books', label: 'Books', icon: 'book', count: 2 },
        { id: 'pdf-papers', label: 'Papers', icon: 'file', count: 8 },
      ]
    }
  ];

  let activeItem = $state('all');

  function getIcon(name: string) {
    switch (name) {
      case 'inbox':
        return '<svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="22 12 16 12 14 15 10 15 8 12 2 12"/><path d="M5.45 5.11 2 12v6a2 2 0 0 0 2 2h16a2 2 0 0 0 2-2v-6l-3.45-6.89A2 2 0 0 0 16.76 4H7.24a2 2 0 0 0-1.79 1.11z"/></svg>';
      case 'circle':
        return '<svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"/></svg>';
      case 'rss':
        return '<svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M4 11a9 9 0 0 1 9 9"/><path d="M4 4a16 16 0 0 1 16 16"/><circle cx="5" cy="19" r="1"/></svg>';
      case 'mail':
        return '<svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect width="20" height="16" x="2" y="4" rx="2"/><path d="m22 7-8.97 5.7a1.94 1.94 0 0 1-2.06 0L2 7"/></svg>';
      case 'book':
        return '<svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M4 19.5v-15A2.5 2.5 0 0 1 6.5 2H20v20H6.5a2.5 2.5 0 0 1 0-5H20"/></svg>';
      case 'file':
        return '<svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M14.5 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V7.5L14.5 2z"/><polyline points="14 2 14 8 20 8"/></svg>';
      default:
        return '';
    }
  }
</script>

<nav class="sidebar-nav">
  <div class="nav-section">
    {#each navItems as item}
      <button
        class="nav-item"
        class:active={activeItem === item.id}
        onclick={() => activeItem = item.id}
      >
        <span class="nav-icon" aria-hidden="true">
          {@html getIcon(item.icon)}
        </span>
        <span class="nav-label">{item.label}</span>
        {#if item.count > 0}
          <span class="nav-count">{item.count}</span>
        {/if}
      </button>
    {/each}
  </div>

  {#each feedSections as section}
    <div class="nav-section">
      <div class="section-header">
        <span class="section-label">{section.label}</span>
        <button class="section-add" title="Add {section.label}">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M12 5v14M5 12h14"/>
          </svg>
        </button>
      </div>
      {#each section.items as item}
        <button
          class="nav-item"
          class:active={activeItem === item.id}
          onclick={() => activeItem = item.id}
        >
          <span class="nav-icon" aria-hidden="true">
            {@html getIcon(item.icon)}
          </span>
          <span class="nav-label">{item.label}</span>
          {#if item.count > 0}
            <span class="nav-count">{item.count}</span>
          {/if}
        </button>
      {/each}
    </div>
  {/each}

  <div class="sidebar-footer">
    <button class="add-feed-btn">
      <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
        <path d="M12 5v14M5 12h14"/>
      </svg>
      Add Feed
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

  .sidebar-footer {
    margin-top: auto;
    padding-top: 8px;
    border-top: 1px solid var(--border-light);
  }

  .add-feed-btn {
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

  .add-feed-btn:hover {
    background: var(--bg-hover);
    border-color: var(--accent);
    color: var(--text-primary);
  }
</style>
