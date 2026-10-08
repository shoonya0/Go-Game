// Fast pre-commit checks: gofmt on staged .go files, then go vet. Full checks
// (tests, builds, lint) run in the Claude Stop hook and in CI (scripts/verify.mjs).

import { spawnSync } from 'node:child_process';

const run = (cmd, args) => spawnSync(cmd, args, { encoding: 'utf8', windowsHide: true });

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
