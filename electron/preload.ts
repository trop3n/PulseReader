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
  }
};

contextBridge.exposeInMainWorld('electronAPI', api);

export type ElectronAPI = typeof api;
