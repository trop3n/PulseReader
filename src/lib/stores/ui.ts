import { writable } from 'svelte/store';

interface UIState {
  showAddFeedModal: boolean;
  selectedFeedId: number | null;
  viewMode: 'all' | 'unread';
}

const initialState: UIState = {
  showAddFeedModal: false,
  selectedFeedId: null,
  viewMode: 'all'
};

function createUIStore() {
  const { subscribe, set, update } = writable<UIState>(initialState);

  return {
    subscribe,
    showAddFeedModal: () => update(s => ({ ...s, showAddFeedModal: true })),
    hideAddFeedModal: () => update(s => ({ ...s, showAddFeedModal: false })),
    selectFeed: (feedId: number | null) => update(s => ({ ...s, selectedFeedId: feedId })),
    setViewMode: (mode: 'all' | 'unread') => update(s => ({ ...s, viewMode: mode })),
    reset: () => set(initialState)
  };
}

export const ui = createUIStore();
