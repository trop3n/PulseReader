import * as esbuild from 'esbuild';
import { mkdirSync, existsSync } from 'fs';
import { join, dirname } from 'path';
import { fileURLToPath } from 'url';

const __dirname = dirname(fileURLToPath(import.meta.url));
const outDir = join(__dirname, '..', 'electron', 'out');

if (!existsSync(outDir)) {
  mkdirSync(outDir, { recursive: true });
}

await esbuild.build({
  entryPoints: [
    join(__dirname, '..', 'electron', 'main.ts'),
    join(__dirname, '..', 'electron', 'preload.ts')
  ],
  outdir: outDir,
  bundle: true,
  platform: 'node',
  target: 'node22',
  format: 'esm',
  external: ['electron', 'better-sqlite3'],
  sourcemap: true
});

console.log('Electron files compiled successfully');
