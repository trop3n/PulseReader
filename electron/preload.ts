import { contextBridge, ipcRenderer } from 'electron';

const api = {
  db: {
    query: (sql: string, params: unknown[] = []) => 
      ipcRenderer.invoke('db:query', sql, params),
    get: (sql: string, params: unknown[] = []) => 
      ipcRenderer.invoke('db:get', sql, params)
  },
  settings: {
    get: (key: string) => ipcRenderer.invoke('settings:get', key),
    set: (key: string, value: unknown) => ipcRenderer.invoke('settings:set', key, value)
  },
  app: {
    getPath: (name: 'home' | 'appData' | 'userData' | 'documents' | 'desktop') => 
      ipcRenderer.invoke('app:getPath', name)
  },
  dialog: {
    openFile: () => ipcRenderer.invoke('dialog:openFile'),
    openDirectory: () => ipcRenderer.invoke('dialog:openDirectory')
  },
  feeds: {
    list: () => ipcRenderer.invoke('feeds:list'),
    add: (name: string, url: string) => ipcRenderer.invoke('feeds:add', { name, url }),
    remove: (feedId: number) => ipcRenderer.invoke('feeds:remove', feedId),
    fetch: (feedId: number) => ipcRenderer.invoke('feeds:fetch', feedId),
    fetchAll: () => ipcRenderer.invoke('feeds:fetchAll')
  },
  articles: {
    list: (options?: { feedId?: number; unreadOnly?: boolean; limit?: number; offset?: number }) => 
      ipcRenderer.invoke('articles:list', options),
    get: (articleId: number) => ipcRenderer.invoke('articles:get', articleId),
    markRead: (articleId: number) => ipcRenderer.invoke('articles:markRead', articleId),
    markAllRead: (feedId?: number) => ipcRenderer.invoke('articles:markAllRead', feedId),
    unreadCount: (feedId?: number) => ipcRenderer.invoke('articles:unreadCount', feedId)
  }
};

contextBridge.exposeInMainWorld('electronAPI', api);

export type ElectronAPI = typeof api;
