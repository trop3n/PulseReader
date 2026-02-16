import { app, BrowserWindow, ipcMain } from 'electron';
import { fileURLToPath } from 'url';
import { dirname, join } from 'path';
import Database from 'better-sqlite3';
import { fetchFeed } from './services/rss.js';

const __filename = fileURLToPath(import.meta.url);
const __dirname = dirname(__filename);

let mainWindow: BrowserWindow | null = null;
let db: Database.Database | null = null;

const isDev = process.env.NODE_ENV === 'development' || !app.isPackaged;

function createWindow() {
  mainWindow = new BrowserWindow({
    width: 1400,
    height: 900,
    minWidth: 1000,
    minHeight: 600,
    backgroundColor: '#1e1e1e',
    titleBarStyle: 'hiddenInset',
    trafficLightPosition: { x: 15, y: 15 },
    webPreferences: {
      preload: join(__dirname, 'preload.js'),
      nodeIntegration: false,
      contextIsolation: true
    }
  });

  if (isDev) {
    mainWindow.loadURL('http://localhost:5173');
    mainWindow.webContents.openDevTools();
  } else {
    mainWindow.loadFile(join(__dirname, '../build/index.html'));
  }

  mainWindow.on('closed', () => {
    mainWindow = null;
  });
}

function initDatabase() {
  const userDataPath = app.getPath('userData');
  const dbPath = join(userDataPath, 'pulsereader.db');
  
  db = new Database(dbPath);
  db.pragma('journal_mode = WAL');
  
  db.exec(`
    CREATE TABLE IF NOT EXISTS feeds (
      id INTEGER PRIMARY KEY AUTOINCREMENT,
      type TEXT NOT NULL CHECK(type IN ('rss', 'email')),
      name TEXT NOT NULL,
      url TEXT,
      config TEXT,
      last_fetched TEXT,
      created_at TEXT DEFAULT CURRENT_TIMESTAMP
    );

    CREATE TABLE IF NOT EXISTS articles (
      id INTEGER PRIMARY KEY AUTOINCREMENT,
      feed_id INTEGER,
      title TEXT NOT NULL,
      content TEXT,
      url TEXT,
      author TEXT,
      published_at TEXT,
      is_read INTEGER DEFAULT 0,
      created_at TEXT DEFAULT CURRENT_TIMESTAMP,
      FOREIGN KEY (feed_id) REFERENCES feeds(id),
      UNIQUE(feed_id, url) ON CONFLICT IGNORE
    );

    CREATE TABLE IF NOT EXISTS pdfs (
      id INTEGER PRIMARY KEY AUTOINCREMENT,
      filepath TEXT NOT NULL UNIQUE,
      title TEXT,
      current_page INTEGER DEFAULT 1,
      total_pages INTEGER,
      created_at TEXT DEFAULT CURRENT_TIMESTAMP
    );

    CREATE TABLE IF NOT EXISTS highlights (
      id INTEGER PRIMARY KEY AUTOINCREMENT,
      source_type TEXT NOT NULL CHECK(source_type IN ('article', 'pdf')),
      source_id INTEGER NOT NULL,
      text TEXT NOT NULL,
      note TEXT,
      position TEXT,
      color TEXT DEFAULT '#7c3aed',
      created_at TEXT DEFAULT CURRENT_TIMESTAMP
    );

    CREATE TABLE IF NOT EXISTS vaults (
      id INTEGER PRIMARY KEY AUTOINCREMENT,
      path TEXT NOT NULL UNIQUE,
      name TEXT,
      is_default INTEGER DEFAULT 0,
      created_at TEXT DEFAULT CURRENT_TIMESTAMP
    );

    CREATE TABLE IF NOT EXISTS settings (
      key TEXT PRIMARY KEY,
      value TEXT
    );
  `);
}

app.whenReady().then(() => {
  initDatabase();
  createWindow();

  app.on('activate', () => {
    if (BrowserWindow.getAllWindows().length === 0) {
      createWindow();
    }
  });
});

app.on('window-all-closed', () => {
  if (process.platform !== 'darwin') {
    app.quit();
  }
});

app.on('before-quit', () => {
  if (db) {
    db.close();
  }
});

ipcMain.handle('db:query', (_event, sql: string, params: unknown[] = []) => {
  if (!db) throw new Error('Database not initialized');
  
  if (sql.trim().toUpperCase().startsWith('SELECT')) {
    return db.prepare(sql).all(...params);
  }
  return db.prepare(sql).run(...params);
});

ipcMain.handle('db:get', (_event, sql: string, params: unknown[] = []) => {
  if (!db) throw new Error('Database not initialized');
  return db.prepare(sql).get(...params);
});

ipcMain.handle('settings:get', (_event, key: string) => {
  if (!db) throw new Error('Database not initialized');
  const row = db.prepare('SELECT value FROM settings WHERE key = ?').get(key) as { value: string } | undefined;
  return row?.value ? JSON.parse(row.value) : null;
});

ipcMain.handle('settings:set', (_event, key: string, value: unknown) => {
  if (!db) throw new Error('Database not initialized');
  db.prepare('INSERT OR REPLACE INTO settings (key, value) VALUES (?, ?)').run(key, JSON.stringify(value));
  return true;
});

ipcMain.handle('app:getPath', (_event, name: string) => {
  return app.getPath(name as Parameters<typeof app.getPath>[0]);
});

ipcMain.handle('dialog:openFile', async () => {
  const { dialog } = await import('electron');
  const result = await dialog.showOpenDialog(mainWindow!, {
    properties: ['openFile'],
    filters: [
      { name: 'PDF Files', extensions: ['pdf'] }
    ]
  });
  return result.filePaths;
});

ipcMain.handle('dialog:openDirectory', async () => {
  const { dialog } = await import('electron');
  const result = await dialog.showOpenDialog(mainWindow!, {
    properties: ['openDirectory']
  });
  return result.filePaths;
});

ipcMain.handle('feeds:list', () => {
  if (!db) throw new Error('Database not initialized');
  return db.prepare(`
    SELECT f.*, 
      (SELECT COUNT(*) FROM articles a WHERE a.feed_id = f.id) as total_count,
      (SELECT COUNT(*) FROM articles a WHERE a.feed_id = f.id AND a.is_read = 0) as unread_count
    FROM feeds f
    ORDER BY f.name
  `).all();
});

ipcMain.handle('feeds:add', async (_event, { name, url }: { name: string; url: string }) => {
  if (!db) throw new Error('Database not initialized');
  
  const existing = db.prepare('SELECT id FROM feeds WHERE url = ?').get(url);
  if (existing) {
    throw new Error('Feed with this URL already exists');
  }

  const result = db.prepare('INSERT INTO feeds (type, name, url) VALUES (?, ?, ?)').run('rss', name, url);
  return { id: result.lastInsertRowid, name, url, type: 'rss' };
});

ipcMain.handle('feeds:remove', (_event, feedId: number) => {
  if (!db) throw new Error('Database not initialized');
  
  db.prepare('DELETE FROM articles WHERE feed_id = ?').run(feedId);
  db.prepare('DELETE FROM feeds WHERE id = ?').run(feedId);
  return true;
});

async function doFetchFeed(feedId: number): Promise<{ success: boolean; count: number }> {
  if (!db) throw new Error('Database not initialized');
  
  const feed = db.prepare('SELECT * FROM feeds WHERE id = ?').get(feedId) as { id: number; url: string; name: string } | undefined;
  if (!feed) {
    throw new Error('Feed not found');
  }

  const feedData = await fetchFeed(feed.url);
  
  const insertArticle = db.prepare(`
    INSERT OR IGNORE INTO articles (feed_id, title, content, url, author, published_at)
    VALUES (?, ?, ?, ?, ?, ?)
  `);

  const insertMany = db.transaction((items: typeof feedData.items) => {
    for (const item of items) {
      insertArticle.run(
        feedId,
        item.title,
        item.content || item.contentSnippet,
        item.link,
        item.creator,
        item.isoDate || item.pubDate
      );
    }
  });

  insertMany(feedData.items);

  db.prepare('UPDATE feeds SET name = ?, last_fetched = CURRENT_TIMESTAMP WHERE id = ?').run(
    feedData.title || feed.name,
    feedId
  );

  return { success: true, count: feedData.items.length };
}

ipcMain.handle('feeds:fetch', async (_event, feedId: number) => {
  try {
    return await doFetchFeed(feedId);
  } catch (error) {
    throw new Error(`Failed to fetch feed: ${(error as Error).message}`);
  }
});

ipcMain.handle('feeds:fetchAll', async () => {
  if (!db) throw new Error('Database not initialized');
  
  const feedsList = db.prepare('SELECT id FROM feeds WHERE type = ?').all('rss') as { id: number }[];
  const results = [];

  for (const feed of feedsList) {
    try {
      const result = await doFetchFeed(feed.id);
      results.push({ feedId: feed.id, ...result });
    } catch (error) {
      results.push({ feedId: feed.id, success: false, error: (error as Error).message });
    }
  }

  return results;
});

ipcMain.handle('articles:list', (_event, options: { feedId?: number; unreadOnly?: boolean; limit?: number; offset?: number } = {}) => {
  if (!db) throw new Error('Database not initialized');
  
  const { feedId, unreadOnly, limit = 50, offset = 0 } = options;
  
  let sql = `
    SELECT a.*, f.name as feed_name, f.type as feed_type
    FROM articles a
    JOIN feeds f ON a.feed_id = f.id
    WHERE 1=1
  `;
  const params: (string | number)[] = [];

  if (feedId) {
    sql += ' AND a.feed_id = ?';
    params.push(feedId);
  }

  if (unreadOnly) {
    sql += ' AND a.is_read = 0';
  }

  sql += ' ORDER BY a.published_at DESC, a.created_at DESC LIMIT ? OFFSET ?';
  params.push(limit, offset);

  return db.prepare(sql).all(...params);
});

ipcMain.handle('articles:get', (_event, articleId: number) => {
  if (!db) throw new Error('Database not initialized');
  
  return db.prepare(`
    SELECT a.*, f.name as feed_name
    FROM articles a
    JOIN feeds f ON a.feed_id = f.id
    WHERE a.id = ?
  `).get(articleId);
});

ipcMain.handle('articles:markRead', (_event, articleId: number) => {
  if (!db) throw new Error('Database not initialized');
  
  db.prepare('UPDATE articles SET is_read = 1 WHERE id = ?').run(articleId);
  return true;
});

ipcMain.handle('articles:markAllRead', (_event, feedId?: number) => {
  if (!db) throw new Error('Database not initialized');
  
  if (feedId) {
    db.prepare('UPDATE articles SET is_read = 1 WHERE feed_id = ?').run(feedId);
  } else {
    db.prepare('UPDATE articles SET is_read = 1').run();
  }
  return true;
});

ipcMain.handle('articles:unreadCount', (_event, feedId?: number) => {
  if (!db) throw new Error('Database not initialized');
  
  if (feedId) {
    const row = db.prepare('SELECT COUNT(*) as count FROM articles WHERE feed_id = ? AND is_read = 0').get(feedId) as { count: number };
    return row.count;
  }
  
  const row = db.prepare('SELECT COUNT(*) as count FROM articles WHERE is_read = 0').get() as { count: number };
  return row.count;
});
