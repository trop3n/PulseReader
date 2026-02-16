<script lang="ts">
  import ArticleList from '$lib/components/articles/ArticleList.svelte';
  import Reader from '$lib/components/reader/Reader.svelte';
  
  let selectedArticle = $state<string | null>(null);
</script>

<div class="dashboard">
  <div class="list-panel">
    <div class="panel-header">
      <h2>All Items</h2>
      <div class="panel-actions">
        <button class="icon-btn" title="Refresh">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M21 12a9 9 0 1 1-9-9c2.52 0 4.93 1 6.74 2.74L21 8"/>
            <path d="M21 3v5h-5"/>
          </svg>
        </button>
        <button class="icon-btn" title="Mark all read">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"/>
            <polyline points="22 4 12 14.01 9 11.01"/>
          </svg>
        </button>
      </div>
    </div>
    <ArticleList onSelect={(id) => selectedArticle = id} />
  </div>
  
  <div class="reader-panel">
    <Reader articleId={selectedArticle} />
  </div>
</div>

<style>
  .dashboard {
    display: flex;
    flex: 1;
    overflow: hidden;
  }

  .list-panel {
    width: var(--list-width);
    min-width: var(--list-width);
    display: flex;
    flex-direction: column;
    background: var(--bg-primary);
    border-right: 1px solid var(--border-light);
  }

  .panel-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 12px 16px;
    border-bottom: 1px solid var(--border-light);
  }

  .panel-header h2 {
    font-size: 14px;
    font-weight: 600;
    color: var(--text-primary);
  }

  .panel-actions {
    display: flex;
    gap: 4px;
  }

  .icon-btn {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 28px;
    height: 28px;
    border-radius: var(--radius-sm);
    color: var(--text-secondary);
    transition: all var(--transition-fast);
  }

  .icon-btn:hover {
    background: var(--bg-hover);
    color: var(--text-primary);
  }

  .reader-panel {
    flex: 1;
    display: flex;
    flex-direction: column;
    overflow: hidden;
  }
</style>
