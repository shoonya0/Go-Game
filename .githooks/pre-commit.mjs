// Fast pre-commit checks: docs-check (errors only), then gofmt on staged .go files and
// go vet. Full checks (tests, builds, lint) run in the Claude Stop hook and in CI
// (scripts/verify.mjs). Git runs hooks from the top of the working tree.

import { spawnSync } from 'node:child_process';
import { existsSync } from 'node:fs';

const run = (cmd, args) => spawnSync(cmd, args, { encoding: 'utf8', windowsHide: true });

// Docs tree: broken links, orphan feature indexes, missing state.md block the commit.
// Budget and stale-state warnings are printed only; CI fails on them for pull requests.
if (existsSync('scripts/docs-check.mjs')) {
  const docs = spawnSync(process.execPath, ['scripts/docs-check.mjs'], { stdio: 'inherit', windowsHide: true });
  if (docs.status !== 0) {
    process.stderr.write('pre-commit: docs-check found errors (fix them, or run /feature-docs check)\n');
    process.exit(docs.status || 1);
  }
}

const staged = run('git', ['diff', '--cached', '--name-only', '--diff-filter=ACMR'])
  .stdout.split(/\r?\n/)
  .filter((f) => f.endsWith('.go'));

if (staged.length === 0) process.exit(0);

const fmt = run('gofmt', ['-l', ...staged]);
if (fmt.error) {
  process.stderr.write(`pre-commit: gofmt not found (${fmt.error.message})\n`);
  process.exit(1);
}
if (fmt.stdout.trim()) {
  process.stderr.write(`pre-commit: unformatted Go files — run gofmt -w on:\n${fmt.stdout}`);
  process.exit(1);
}

const vet = spawnSync('go', ['vet', './...'], { stdio: 'inherit', windowsHide: true });
if (vet.status !== 0) {
  process.stderr.write('pre-commit: go vet failed\n');
  process.exit(vet.status || 1);
}
