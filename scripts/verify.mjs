#!/usr/bin/env node
// scripts/verify.mjs — one cross-platform "verify" entry point for this Go project.
// Runs the fast checks in order and stops at the first failure; `--full` adds the
// native + WebAssembly builds and a full (advisory) lint report.
//
//   node scripts/verify.mjs          # fast: what the Stop hook runs (~30s)
//   node scripts/verify.mjs --full   # fast + slow: what CI runs
//
// .claude/settings.json sets CLAUDE_VERIFY_CMD to "node scripts/verify.mjs".
// Commands are constant argv arrays from this file, run without a shell.

import { spawnSync } from 'node:child_process';
import { existsSync } from 'node:fs';
import { tmpdir } from 'node:os';
import { join } from 'node:path';

// golangci-lint is usually installed with `go install` into $(go env GOPATH)/bin,
// which is often not on PATH. Prefer PATH, fall back to GOPATH/bin.
function resolveGoTool(name) {
  const exe = process.platform === 'win32' ? `${name}.exe` : name;
  const onPath = spawnSync(name, ['--version'], { stdio: 'ignore', windowsHide: true });
  if (!onPath.error) return name;
  const gopath = spawnSync('go', ['env', 'GOPATH'], { encoding: 'utf8', windowsHide: true }).stdout?.trim();
  const candidate = gopath ? join(gopath, 'bin', exe) : '';
  return candidate && existsSync(candidate) ? candidate : null;
}

const lint = resolveGoTool('golangci-lint');
// CI checkouts have no local main branch; the workflow sets VERIFY_BASE=origin/main.
const base = process.env.VERIFY_BASE || 'main';
const lintMissing = 'golangci-lint not found: go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest';

const FAST = [
  { name: 'format', cmd: ['gofmt', '-l', '.'], failOnOutput: 'unformatted files (run gofmt -w <file>):' },
  { name: 'vet', cmd: ['go', 'vet', './...'] },
  // Only issues introduced since main block; the baseline issues are listed in docs/current-state.md.
  { name: 'lint (new issues)', cmd: lint ? [lint, 'run', `--new-from-merge-base=${base}`, './...'] : null, missing: lintMissing },
  { name: 'unit tests', cmd: ['go', 'test', './...'] },
];

const SLOW = [
  { name: 'build (native)', cmd: ['go', 'build', './...'] },
  {
    name: 'build (wasm)',
    cmd: ['go', 'build', '-ldflags=-s -w', '-o', join(tmpdir(), 'go-game-verify.wasm'), './cmd'],
    env: { GOOS: 'js', GOARCH: 'wasm' },
  },
  { name: 'lint (full, advisory)', cmd: lint ? [lint, 'run', '--issues-exit-code=0', './...'] : null, missing: lintMissing },
];

const steps = process.argv.includes('--full') ? [...FAST, ...SLOW] : FAST;
const summary = [];

function fail(step, code, why) {
  summary.push(`✘ ${step.name}${why ? ` — ${why}` : ''}`);
  process.stdout.write(`\n${summary.join('\n')}\nverify FAILED at "${step.name}"\n`);
  process.exit(code || 1);
}

for (const step of steps) {
  if (!step.cmd) fail(step, 1, step.missing);
  const started = Date.now();
  process.stdout.write(`\n▶ ${step.name}: ${step.cmd.join(' ')}\n`);
  const r = spawnSync(step.cmd[0], step.cmd.slice(1), {
    stdio: step.failOnOutput ? ['ignore', 'pipe', 'inherit'] : 'inherit',
    encoding: 'utf8',
    env: { ...process.env, ...step.env },
    windowsHide: true,
  });
  const seconds = ((Date.now() - started) / 1000).toFixed(1);
  if (r.error) fail(step, 1, r.error.message);
  if (step.failOnOutput && r.stdout.trim()) {
    process.stdout.write(`${step.failOnOutput}\n${r.stdout}`);
    fail(step, 1, `${seconds}s`);
  }
  if (r.status !== 0) fail(step, r.status, `${seconds}s`);
  summary.push(`✔ ${step.name} (${seconds}s)`);
}

process.stdout.write(`\n${summary.join('\n')}\nverify passed (${steps.length} steps)\n`);
