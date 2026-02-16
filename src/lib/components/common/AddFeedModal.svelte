<script lang="ts">
  import { feeds } from '$lib/stores/feeds';

  interface Props {
    open: boolean;
    onClose: () => void;
  }

  let { open, onClose }: Props = $props();

  let url = $state('');
  let name = $state('');
  let loading = $state(false);
  let error = $state<string | null>(null);

  async function handleSubmit(e: Event) {
    e.preventDefault();
    
    if (!url.trim()) {
      error = 'Please enter a feed URL';
      return;
    }

    loading = true;
    error = null;

    try {
      const feedName = name.trim() || extractFeedName(url);
      const feed = await feeds.add(feedName, url.trim());
      
      await feeds.fetch(feed.id);
      
      url = '';
      name = '';
      onClose();
    } catch (e) {
      error = (e as Error).message;
    } finally {
      loading = false;
    }
  }

  function extractFeedName(url: string): string {
    try {
      const urlObj = new URL(url);
      return urlObj.hostname.replace('www.', '');
    } catch {
      return 'New Feed';
    }
  }

  function handleKeydown(e: KeyboardEvent) {
    if (e.key === 'Escape') {
      onClose();
    }
  }

  function reset() {
    url = '';
    name = '';
    error = null;
    loading = false;
  }
</script>

<svelte:window on:keydown={handleKeydown} />

{#if open}
  <!-- svelte-ignore a11y_no_static_element_interactions a11y_click_events_have_key_events -->
  <div class="modal-overlay" onclick={onClose} onkeydown={handleKeydown} role="button" tabindex={-1}>
    <!-- svelte-ignore a11y_no_static_element_interactions a11y_click_events_have_key_events -->
    <div class="modal" onclick={(e) => e.stopPropagation()} onkeydown={(e) => e.stopPropagation()} role="dialog" aria-modal="true" tabindex={-1}>
      <div class="modal-header">
        <h2>Add Feed</h2>
        <button class="close-btn" onclick={onClose} aria-label="Close">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M18 6 6 18M6 6l12 12"/>
          </svg>
        </button>
      </div>

      <form onsubmit={handleSubmit}>
        <div class="form-group">
          <label for="feed-url">Feed URL</label>
          <input
            id="feed-url"
            type="url"
            bind:value={url}
            placeholder="https://example.com/feed.xml"
            disabled={loading}
          />
        </div>

        <div class="form-group">
          <label for="feed-name">Name (optional)</label>
          <input
            id="feed-name"
            type="text"
            bind:value={name}
            placeholder="Auto-detected from feed"
            disabled={loading}
          />
        </div>

        {#if error}
          <div class="error-message">{error}</div>
        {/if}

        <div class="modal-actions">
          <button type="button" class="btn-secondary" onclick={onClose} disabled={loading}>
            Cancel
          </button>
          <button type="submit" class="btn-primary" disabled={loading}>
            {loading ? 'Adding...' : 'Add Feed'}
          </button>
        </div>
      </form>
    </div>
  </div>
{/if}

<style>
  .modal-overlay {
    position: fixed;
    inset: 0;
    background: rgba(0, 0, 0, 0.6);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 1000;
    animation: fadeIn 150ms ease;
  }

  @keyframes fadeIn {
    from { opacity: 0; }
    to { opacity: 1; }
  }

  .modal {
    background: var(--bg-secondary);
    border: 1px solid var(--border);
    border-radius: var(--radius-lg);
    width: 100%;
    max-width: 420px;
    margin: 16px;
    box-shadow: var(--shadow-lg);
    animation: slideIn 150ms ease;
  }

  @keyframes slideIn {
    from { 
      opacity: 0;
      transform: translateY(-10px);
    }
    to { 
      opacity: 1;
      transform: translateY(0);
    }
  }

  .modal-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 16px 20px;
    border-bottom: 1px solid var(--border-light);
  }

  .modal-header h2 {
    font-size: 16px;
    font-weight: 600;
    color: var(--text-primary);
  }

  .close-btn {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 28px;
    height: 28px;
    border-radius: var(--radius-sm);
    color: var(--text-secondary);
    transition: all var(--transition-fast);
  }

  .close-btn:hover {
    background: var(--bg-hover);
    color: var(--text-primary);
  }

  form {
    padding: 20px;
  }

  .form-group {
    margin-bottom: 16px;
  }

  .form-group label {
    display: block;
    font-size: 12px;
    font-weight: 500;
    color: var(--text-secondary);
    margin-bottom: 6px;
  }

  .form-group input {
    width: 100%;
    padding: 10px 12px;
    background: var(--bg-primary);
    border: 1px solid var(--border);
    border-radius: var(--radius-sm);
    color: var(--text-primary);
    font-size: 14px;
    transition: border-color var(--transition-fast);
  }

  .form-group input:focus {
    outline: none;
    border-color: var(--accent);
  }

  .form-group input::placeholder {
    color: var(--text-muted);
  }

  .error-message {
    padding: 10px 12px;
    margin-bottom: 16px;
    background: rgba(248, 113, 113, 0.1);
    border: 1px solid var(--error);
    border-radius: var(--radius-sm);
    font-size: 13px;
    color: var(--error);
  }

  .modal-actions {
    display: flex;
    justify-content: flex-end;
    gap: 8px;
    margin-top: 8px;
  }

  .btn-secondary,
  .btn-primary {
    padding: 8px 16px;
    border-radius: var(--radius-sm);
    font-size: 13px;
    font-weight: 500;
    transition: all var(--transition-fast);
  }

  .btn-secondary {
    background: var(--bg-tertiary);
    border: 1px solid var(--border);
    color: var(--text-secondary);
  }

  .btn-secondary:hover:not(:disabled) {
    background: var(--bg-hover);
    color: var(--text-primary);
  }

  .btn-primary {
    background: var(--accent);
    border: 1px solid var(--accent);
    color: white;
  }

  .btn-primary:hover:not(:disabled) {
    background: var(--accent-hover);
  }

  .btn-secondary:disabled,
  .btn-primary:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }
</style>
