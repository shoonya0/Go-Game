// PostToolUse hook (matcher: Edit|Write)
//
// Keeps the code-review-graph knowledge graph current after every edit, so graph
// queries (callers_of, tests_for, impact radius) never answer from stale data.
// Silent no-op when code-review-graph isn't installed or the project isn't a git repo.
//
// Spawned without a shell. On Windows a pip/pipx/uv-installed console script is an
// .exe on PATH, which spawnSync resolves directly; if it isn't found we exit 0.

import { spawnSync } from 'node:child_process';

function readStdin() {
  return new Promise((resolve) => {
    let raw = '';
    process.stdin.setEncoding('utf8');
    process.stdin.on('data', (chunk) => (raw += chunk));
    process.stdin.on('end', () => resolve(raw));
    process.stdin.on('error', () => resolve(raw));
  });
}

// Read the hook payload first so the runner never blocks on a full stdin pipe
// (Write payloads include the whole file content).
await readStdin();

const projectDir = process.env.CLAUDE_PROJECT_DIR || process.cwd();

spawnSync('code-review-graph', ['update', '--skip-flows', '--repo', projectDir], {
  cwd: projectDir,
  encoding: 'utf8',
  timeout: 25_000,
  windowsHide: true,
});

// Not installed (ENOENT), timeouts and non-git folders are all non-fatal: a stale
// graph is better than an interrupted edit loop.
process.exit(0);
