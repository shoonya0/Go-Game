#!/usr/bin/env node
// check-tools.mjs — READ-ONLY detection of the tools the AI harness needs.
//
// Usage:  node check-tools.mjs            → human-readable table
//         node check-tools.mjs --json     → machine-readable JSON (for the agent)
//
// Never installs anything. Finds executables by scanning PATH (+ PATHEXT on Windows)
// with fs only, then runs fixed, constant version commands through the platform
// shell (needed for Windows .cmd shims such as npm.cmd). No user input reaches a shell.

import { spawnSync } from 'node:child_process';
import { existsSync, statSync } from 'node:fs';
import path from 'node:path';

const isWin = process.platform === 'win32';

/** Find an executable on PATH without spawning anything. */
function which(name) {
  const dirs = (process.env.PATH || process.env.Path || '').split(path.delimiter).filter(Boolean);
  const exts = isWin
    ? (process.env.PATHEXT || '.COM;.EXE;.BAT;.CMD').split(';').filter(Boolean).map((e) => e.toLowerCase())
    : [''];
  for (const dir of dirs) {
    for (const ext of isWin ? ['', ...exts] : exts) {
      const candidate = path.join(dir, name + ext);
      try {
        if (existsSync(candidate) && statSync(candidate).isFile()) {
          // On Windows a bare name without extension is not runnable; require an extension.
          if (isWin && !path.extname(candidate)) continue;
          return candidate;
        }
      } catch {
        /* unreadable PATH entry: ignore */
      }
    }
  }
  return null;
}

/** Run a constant command string; return { ok, out }. */
function run(command, timeout = 20_000) {
  const r = spawnSync(command, { shell: true, encoding: 'utf8', timeout, windowsHide: true });
  const out = `${r.stdout || ''}${r.stderr || ''}`.trim().split(/\r?\n/)[0] || '';
  return { ok: r.status === 0 && !r.error, out };
}

const TOOLS = [
  { id: 'git', tier: 0, names: ['git'], version: 'git --version', why: 'version control; Git Bash on Windows' },
  { id: 'node', tier: 0, names: ['node'], version: 'node --version', why: 'runs the hook scripts, agent-browser, npx tools', minMajor: 20 },
  { id: 'npm', tier: 0, names: ['npm'], version: 'npm --version', why: 'installs JS tools' },
  { id: 'gh', tier: 0, names: ['gh'], version: 'gh --version', why: 'PRs, issues, CI logs', auth: 'gh auth status' },
  { id: 'claude', tier: 0, names: ['claude'], version: 'claude --version', why: 'the coding agent CLI (plugin/mcp checks)' },
  { id: 'python', tier: 1, names: isWin ? ['python', 'py'] : ['python3', 'python'], version: '{bin} --version', why: 'needed for pip-installed tools (or use uv)' },
  { id: 'uv', tier: 1, names: ['uv'], version: 'uv --version', why: 'preferred installer for Python CLI tools' },
  { id: 'pipx', tier: 1, names: ['pipx'], version: 'pipx --version', why: 'alternative installer for Python CLI tools', optional: true },
  { id: 'code-review-graph', tier: 1, names: ['code-review-graph'], version: null, why: 'code knowledge graph (MCP + CLI)' },
  { id: 'agent-browser', tier: 2, names: ['agent-browser'], version: 'agent-browser --version', why: 'agent drives web/Electron UI over CDP (UI projects only)' },
];

const results = [];
for (const tool of TOOLS) {
  let found = null;
  let bin = null;
  for (const n of tool.names) {
    found = which(n);
    if (found) {
      bin = n;
      break;
    }
  }
  const entry = { id: tool.id, tier: tool.tier, optional: !!tool.optional, why: tool.why, present: !!found, path: found, version: null, working: null };
  if (found && tool.version) {
    const v = run(tool.version.replace('{bin}', bin));
    entry.working = v.ok;
    entry.version = v.out || null;
    // Windows Store "python" alias stub: present on PATH but not a real interpreter.
    if (tool.id === 'python' && found.toLowerCase().includes('windowsapps') && !v.ok) entry.present = false;
  } else if (found) {
    entry.working = true;
  }
  if (tool.minMajor && entry.version) {
    const major = Number.parseInt(String(entry.version).replace(/^v/, ''), 10);
    entry.meetsMinimum = Number.isFinite(major) ? major >= tool.minMajor : null;
    entry.minimum = `>=${tool.minMajor}`;
  }
  if (found && tool.auth) entry.authenticated = run(tool.auth).ok;
  results.push(entry);
}

const report = {
  platform: process.platform,
  arch: process.arch,
  shellHint: isWin ? 'Windows: Claude Code uses Git Bash for the Bash tool when Git for Windows is installed, else PowerShell' : 'POSIX shell',
  tools: results,
  missingRequired: results.filter((r) => r.tier === 0 && (!r.present || r.working === false || r.meetsMinimum === false)).map((r) => r.id),
  missingRecommended: results.filter((r) => r.tier === 1 && !r.optional && !r.present).map((r) => r.id),
};

if (process.argv.includes('--json')) {
  process.stdout.write(`${JSON.stringify(report, null, 2)}\n`);
} else {
  console.log(`platform: ${report.platform} (${report.arch})`);
  for (const r of results) {
    const state = !r.present ? 'MISSING' : r.working === false ? 'BROKEN' : r.meetsMinimum === false ? `TOO OLD (${r.minimum})` : 'ok';
    const extra = r.authenticated === false ? ' [not logged in]' : r.authenticated ? ' [logged in]' : '';
    console.log(`T${r.tier} ${r.id.padEnd(18)} ${state.padEnd(16)} ${r.version || ''}${extra}`);
  }
  console.log(`missing required: ${report.missingRequired.join(', ') || 'none'}`);
  console.log(`missing recommended: ${report.missingRecommended.join(', ') || 'none'}`);
}
