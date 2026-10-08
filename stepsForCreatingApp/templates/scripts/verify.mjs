#!/usr/bin/env node
// scripts/verify.mjs — one cross-platform "verify" entry point for projects that
// don't use npm scripts (Python, Go, Rust, Java, .NET …). Runs the fast checks in
// order and stops at the first failure; `--full` adds the slow checks.
//
//   node scripts/verify.mjs          # fast: what the Stop hook runs
//   node scripts/verify.mjs --full   # fast + slow: what CI runs
//
// Set CLAUDE_VERIFY_CMD to "node scripts/verify.mjs" in .claude/settings.json "env".
// Edit STEPS for your stack (see stepsForCreatingApp/agent/reference/stack-matrix.md).
// Commands are constant strings from this file, run through the platform shell so
// tool shims (.cmd/.bat on Windows) resolve; never build them from untrusted input.

import { spawnSync } from 'node:child_process';

const FAST = [
  { name: 'typecheck', cmd: 'pyright' }, //        e.g. go vet ./...   | cargo check
  { name: 'lint', cmd: 'ruff check .' }, //         e.g. golangci-lint run | cargo clippy -- -D warnings
  { name: 'format', cmd: 'ruff format --check .' }, // e.g. gofmt -l .  | cargo fmt --check
  { name: 'unit tests', cmd: 'pytest -q' }, //      e.g. go test ./...  | cargo test
];

const SLOW = [
  // { name: 'integration tests', cmd: 'pytest -q -m integration' },
  // { name: 'build', cmd: 'python -m build' },
];

const steps = process.argv.includes('--full') ? [...FAST, ...SLOW] : FAST;
const summary = [];

for (const step of steps) {
  const started = Date.now();
  process.stdout.write(`\n▶ ${step.name}: ${step.cmd}\n`);
  const r = spawnSync(step.cmd, { shell: true, stdio: 'inherit', windowsHide: true });
  const seconds = ((Date.now() - started) / 1000).toFixed(1);
  const ok = r.status === 0 && !r.error;
  summary.push(`${ok ? '✔' : '✘'} ${step.name} (${seconds}s)${r.error ? ` — ${r.error.message}` : ''}`);
  if (!ok) {
    process.stdout.write(`\n${summary.join('\n')}\nverify FAILED at "${step.name}"\n`);
    process.exit(r.status || 1);
  }
}

process.stdout.write(`\n${summary.join('\n')}\nverify passed (${steps.length} steps)\n`);
