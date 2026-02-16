<script lang="ts">
  interface Props {
    articleId?: string | null;
  }

  let { articleId }: Props = $props();

  const sampleContent = `
    <h1>Understanding React Server Components</h1>
    <p class="meta">Published by Tech Blogs · 2 hours ago</p>
    <p>React Server Components represent a fundamental shift in how we think about React applications. By moving certain components to the server, we can reduce the amount of JavaScript sent to the client and improve performance.</p>
    <h2>What Are Server Components?</h2>
    <p>Server Components are a new type of React component that renders only on the server. Unlike traditional React components, they don't add to the client-side JavaScript bundle.</p>
    <p>This means you can have rich, interactive applications without the performance penalty of large JavaScript bundles. The key insight is that not all components need to be interactive.</p>
    <h2>Benefits</h2>
    <ul>
      <li><strong>Smaller Bundle Size:</strong> Server Components and their dependencies are not included in the client bundle.</li>
      <li><strong>Direct Backend Access:</strong> Server Components can directly access backend resources like databases and file systems.</li>
      <li><strong>Automatic Code Splitting:</strong> Only interactive parts of your app load client-side code.</li>
    </ul>
    <blockquote>
      "Server Components let you render parts of your application on the server, streaming them to the client. This is a fundamental shift in React's mental model." — React Team
    </blockquote>
    <h2>Getting Started</h2>
    <p>To use Server Components, you'll need React 18 or later and a framework that supports them, such as Next.js 13+.</p>
  `;
</script>

<div class="reader">
  {#if articleId}
    <div class="reader-toolbar">
      <div class="toolbar-actions">
        <button class="toolbar-btn" title="Previous article">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="m15 18-6-6 6-6"/>
          </svg>
        </button>
        <button class="toolbar-btn" title="Next article">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="m9 18 6-6-6-6"/>
          </svg>
        </button>
      </div>
      <div class="toolbar-spacer"></div>
      <button class="toolbar-btn primary" title="Send to Obsidian">
        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/>
          <polyline points="17 8 12 3 7 8"/>
          <line x1="12" x2="12" y1="3" y2="15"/>
        </svg>
        <span>Send to Obsidian</span>
      </button>
    </div>

    <div class="reader-content">
      <article class="article">
        {@html sampleContent}
      </article>

      <aside class="highlights-panel">
        <div class="panel-title">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M12 20h9"/>
            <path d="M16.5 3.5a2.12 2.12 0 0 1 3 3L7 19l-4 1 1-4Z"/>
          </svg>
          Highlights
        </div>
        <div class="highlight-item">
          <p class="highlight-text">"Server Components let you render parts of your application on the server..."</p>
          <textarea class="highlight-note" placeholder="Add a note..."></textarea>
        </div>
      </aside>
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
  }

  .reader-toolbar {
    display: flex;
    align-items: center;
    padding: 8px 16px;
    background: var(--bg-secondary);
    border-bottom: 1px solid var(--border-light);
  }

  .toolbar-actions {
    display: flex;
    gap: 4px;
  }

  .toolbar-spacer {
    flex: 1;
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
    display: flex;
    flex: 1;
    overflow: hidden;
  }

  .article {
    flex: 1;
    padding: 32px 48px;
    overflow-y: auto;
    max-width: 720px;
  }

  .article :global(h1) {
    font-size: 28px;
    font-weight: 600;
    color: var(--text-primary);
    margin-bottom: 8px;
    line-height: 1.3;
  }

  .article :global(.meta) {
    font-size: 13px;
    color: var(--text-muted);
    margin-bottom: 24px;
  }

  .article :global(h2) {
    font-size: 20px;
    font-weight: 600;
    color: var(--text-primary);
    margin-top: 32px;
    margin-bottom: 12px;
  }

  .article :global(p) {
    font-size: 15px;
    line-height: 1.7;
    color: var(--text-primary);
    margin-bottom: 16px;
  }

  .article :global(ul) {
    margin-bottom: 16px;
    padding-left: 24px;
  }

  .article :global(li) {
    font-size: 15px;
    line-height: 1.7;
    color: var(--text-primary);
    margin-bottom: 8px;
  }

  .article :global(blockquote) {
    border-left: 3px solid var(--accent);
    padding-left: 16px;
    margin: 24px 0;
    font-style: italic;
    color: var(--text-secondary);
  }

  .highlights-panel {
    width: 280px;
    background: var(--bg-secondary);
    border-left: 1px solid var(--border-light);
    padding: 16px;
    overflow-y: auto;
  }

  .panel-title {
    display: flex;
    align-items: center;
    gap: 8px;
    font-size: 12px;
    font-weight: 600;
    text-transform: uppercase;
    letter-spacing: 0.5px;
    color: var(--text-muted);
    margin-bottom: 16px;
  }

  .highlight-item {
    background: var(--bg-tertiary);
    border-radius: var(--radius-md);
    padding: 12px;
    margin-bottom: 12px;
  }

  .highlight-text {
    font-size: 13px;
    color: var(--text-primary);
    line-height: 1.5;
    margin-bottom: 8px;
    font-style: italic;
  }

  .highlight-note {
    width: 100%;
    min-height: 60px;
    padding: 8px;
    background: var(--bg-primary);
    border: 1px solid var(--border);
    border-radius: var(--radius-sm);
    font-size: 13px;
    resize: vertical;
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
