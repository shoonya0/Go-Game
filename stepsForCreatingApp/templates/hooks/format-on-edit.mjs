// PostToolUse hook (matcher: Edit|Write)
//
// Formats the file the agent just edited with the PROJECT's own Prettier, so diffs
// stay free of style noise without the agent spending tokens on formatting.
//
// - No-op when Prettier isn't installed in the project or the file type isn't handled.
// - Runs Prettier's JS entry with the current Node (process.execPath) instead of an
//   npx/.cmd shim, so it spawns identically on Windows and macOS with no shell.
// - On a Prettier failure (usually a syntax error) it exits 2: for PostToolUse that
//   shows stderr to the agent (the edit already happened), so it can fix the syntax.
//
// Swap in Biome / ruff / gofmt by changing FORMATTER below.

import { spawnSync } from 'node:child_process';
import { existsSync, readFileSync } from 'node:fs';
import { createRequire } from 'node:module';
import path from 'node:path';

const EXTENSIONS = new Set([
  '.js', '.jsx', '.mjs', '.cjs', '.ts', '.tsx', '.mts', '.cts',
  '.json', '.css', '.scss', '.less', '.html', '.md', '.mdx', '.yaml', '.yml', '.vue', '.svelte',
]);

function readStdin() {
  return new Promise((resolve) => {
    let raw = '';
    process.stdin.setEncoding('utf8');
    process.stdin.on('data', (chunk) => (raw += chunk));
    process.stdin.on('end', () => resolve(raw));
    process.stdin.on('error', () => resolve(raw));
  });
}

function findPrettierBin(projectDir) {
  let pkgJsonPath;
  try {
    pkgJsonPath = createRequire(path.join(projectDir, 'package.json')).resolve('prettier/package.json');
  } catch {
    const direct = path.join(projectDir, 'node_modules', 'prettier', 'package.json');
    if (!existsSync(direct)) return null;
    pkgJsonPath = direct;
  }
  const pkg = JSON.parse(readFileSync(pkgJsonPath, 'utf8'));
  const binRel = typeof pkg.bin === 'string' ? pkg.bin : pkg.bin?.prettier;
  if (!binRel) return null;
  const bin = path.join(path.dirname(pkgJsonPath), binRel);
  return existsSync(bin) ? bin : null;
}

const raw = await readStdin();
let input = {};
try {
  input = JSON.parse(raw || '{}');
} catch {
  process.exit(0);
}

const target = input.tool_input?.file_path;
if (!target) process.exit(0);

const projectDir = process.env.CLAUDE_PROJECT_DIR || input.cwd || process.cwd();
const file = path.resolve(projectDir, target);
if (!EXTENSIONS.has(path.extname(file).toLowerCase()) || !existsSync(file)) process.exit(0);

const FORMATTER = findPrettierBin(projectDir);
if (!FORMATTER) process.exit(0); // project doesn't use Prettier: nothing to do

const result = spawnSync(
  process.execPath,
  [FORMATTER, '--write', '--ignore-unknown', '--log-level', 'warn', file],
  { cwd: projectDir, encoding: 'utf8', timeout: 25_000 },
);

if (result.status !== 0) {
  const output = `${result.stderr ?? ''}${result.stdout ?? ''}`.trim();
  process.stderr.write(`Prettier could not format ${path.relative(projectDir, file)}:\n${output.slice(-2000)}\n`);
  process.exit(2);
}

process.exit(0);
