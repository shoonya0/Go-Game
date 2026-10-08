#!/usr/bin/env node
// inspect-project.mjs — READ-ONLY snapshot of a target folder, used to choose the
// setup path (NEW vs EXISTING) and to avoid overwriting existing AI configuration.
//
// Usage:  node inspect-project.mjs <targetDir> [--json]
//
// Reports: git state, languages (by marker files + source extensions), package
// scripts, CI, pre-commit tooling, existing AI artifacts (CLAUDE.md, AGENTS.md,
// .claude/, .mcp.json …), whether .gitignore hides .claude/, and a suggested mode.
// Never writes anything. Only spawns `git` with fixed arguments.

import { spawnSync } from 'node:child_process';
import { existsSync, readdirSync, readFileSync, statSync } from 'node:fs';
import path from 'node:path';

const args = process.argv.slice(2);
const asJson = args.includes('--json');
const target = path.resolve(args.find((a) => !a.startsWith('--')) || process.cwd());

if (!existsSync(target)) {
  const report = { target, exists: false, suggestedMode: 'NEW', reason: 'target folder does not exist yet' };
  process.stdout.write(asJson ? `${JSON.stringify(report, null, 2)}\n` : `target ${target} does not exist → mode NEW\n`);
  process.exit(0);
}

const has = (rel) => existsSync(path.join(target, rel));
const read = (rel) => {
  try {
    return readFileSync(path.join(target, rel), 'utf8');
  } catch {
    return null;
  }
};
const git = (...gitArgs) => {
  const r = spawnSync('git', ['-C', target, ...gitArgs], { encoding: 'utf8', windowsHide: true });
  return r.status === 0 ? r.stdout.trim() : null;
};

// --- git -------------------------------------------------------------------
const isRepo = git('rev-parse', '--is-inside-work-tree') === 'true';
const gitInfo = isRepo
  ? {
      isRepo: true,
      root: git('rev-parse', '--show-toplevel'),
      branch: git('rev-parse', '--abbrev-ref', 'HEAD'),
      defaultBranch: (git('symbolic-ref', '--short', 'refs/remotes/origin/HEAD') || '').replace(/^origin\//, '') || null,
      remotes: (git('remote') || '').split(/\r?\n/).filter(Boolean),
      dirtyFiles: (git('status', '--porcelain') || '').split(/\r?\n/).filter(Boolean).length,
      commits: Number.parseInt(git('rev-list', '--count', 'HEAD') || '0', 10),
      claudeDirIgnored: git('check-ignore', '-q', '.claude/settings.json') !== null,
      claudeIgnoreRule: git('check-ignore', '-v', '.claude/settings.json'),
    }
  : { isRepo: false };

// --- files -----------------------------------------------------------------
const SKIP = new Set(['.git', 'node_modules', 'dist', 'build', 'out', 'target', 'vendor', '.venv', 'venv', '__pycache__', '.next', '.nuxt', 'coverage', '.turbo', '.cache', 'bin', 'obj', '.gradle', '.idea', '.vscode', 'Pods', '.code-review-graph', 'release']);
const SOURCE_EXT = { '.ts': 'TypeScript', '.tsx': 'TypeScript', '.js': 'JavaScript', '.jsx': 'JavaScript', '.mjs': 'JavaScript', '.cjs': 'JavaScript', '.py': 'Python', '.go': 'Go', '.rs': 'Rust', '.java': 'Java', '.kt': 'Kotlin', '.cs': 'C#', '.swift': 'Swift', '.rb': 'Ruby', '.php': 'PHP', '.c': 'C/C++', '.cc': 'C/C++', '.cpp': 'C/C++', '.h': 'C/C++', '.dart': 'Dart', '.vue': 'Vue', '.svelte': 'Svelte' };
const counts = {};
let sourceFiles = 0;
let testFiles = 0;
let scanned = 0;
const LIMIT = 200_000;
(function walk(dir, depth) {
  if (scanned > LIMIT || depth > 12) return;
  let entries = [];
  try {
    entries = readdirSync(dir, { withFileTypes: true });
  } catch {
    return;
  }
  for (const e of entries) {
    if (scanned > LIMIT) return;
    if (e.isDirectory()) {
      if (!SKIP.has(e.name) && !e.name.startsWith('.')) walk(path.join(dir, e.name), depth + 1);
      continue;
    }
    scanned++;
    const lang = SOURCE_EXT[path.extname(e.name).toLowerCase()];
    if (!lang) continue;
    counts[lang] = (counts[lang] || 0) + 1;
    sourceFiles++;
    if (/(\.|_|-)(test|spec)\.[^.]+$|^test_.*\.py$|_test\.go$/i.test(e.name) || /[\\/](__tests__|tests?)[\\/]/i.test(dir + path.sep)) testFiles++;
  }
})(target, 0);

const MARKERS = { 'package.json': 'Node/JS', 'tsconfig.json': 'TypeScript', 'pyproject.toml': 'Python', 'requirements.txt': 'Python', 'go.mod': 'Go', 'Cargo.toml': 'Rust', 'pom.xml': 'Java (Maven)', 'build.gradle': 'JVM (Gradle)', 'build.gradle.kts': 'JVM (Gradle)', 'Gemfile': 'Ruby', 'composer.json': 'PHP', 'Package.swift': 'Swift', 'pubspec.yaml': 'Dart/Flutter' };
const markers = Object.keys(MARKERS).filter(has).map((m) => `${m} (${MARKERS[m]})`);
const csproj = readdirSync(target).filter((n) => /\.(csproj|sln)$/i.test(n));
if (csproj.length) markers.push(`${csproj[0]} (.NET)`);

// --- package.json details --------------------------------------------------
let pkg = null;
const pkgRaw = read('package.json');
if (pkgRaw) {
  try {
    const p = JSON.parse(pkgRaw);
    const deps = { ...(p.dependencies || {}), ...(p.devDependencies || {}) };
    pkg = {
      name: p.name || null,
      packageManager: p.packageManager || (has('pnpm-lock.yaml') ? 'pnpm' : has('yarn.lock') ? 'yarn' : has('bun.lockb') || has('bun.lock') ? 'bun' : has('package-lock.json') ? 'npm' : null),
      scripts: Object.keys(p.scripts || {}),
      hasScript: Object.fromEntries(['verify', 'test', 'lint', 'typecheck', 'format', 'build', 'dev', 'test:e2e'].map((s) => [s, !!(p.scripts || {})[s]])),
      notable: ['typescript', 'eslint', 'prettier', '@biomejs/biome', 'vitest', 'jest', '@playwright/test', 'cypress', 'electron', 'react', 'next', 'vue', 'svelte', 'husky', 'lefthook', 'dependency-cruiser', 'eslint-plugin-boundaries', 'knip'].filter((d) => d in deps),
    };
  } catch {
    pkg = { error: 'package.json is not valid JSON' };
  }
}

// --- harness / AI artifacts ------------------------------------------------
const list = (rel) => {
  try {
    return readdirSync(path.join(target, rel));
  } catch {
    return [];
  }
};
const ai = {
  'CLAUDE.md': has('CLAUDE.md') || has('.claude/CLAUDE.md'),
  'CLAUDE.local.md': has('CLAUDE.local.md'),
  'AGENTS.md': has('AGENTS.md'),
  '.claude/settings.json': has('.claude/settings.json'),
  '.claude/settings.local.json': has('.claude/settings.local.json'),
  '.claude/skills': list('.claude/skills'),
  '.claude/agents': list('.claude/agents'),
  '.claude/hooks': list('.claude/hooks'),
  '.claude/rules': list('.claude/rules'),
  '.mcp.json': has('.mcp.json'),
  '.cursor / .cursorrules': has('.cursor') || has('.cursorrules'),
  '.github/copilot-instructions.md': has('.github/copilot-instructions.md'),
};
const aiArtifacts = Object.entries(ai).filter(([, v]) => (Array.isArray(v) ? v.length > 0 : v)).map(([k]) => k);
// A graph cache alone is a tool artefact, not evidence of an AI workflow: report it, don't switch mode on it.
const graphCachePresent = has('.code-review-graph');

const harness = {
  ciWorkflows: list('.github/workflows'),
  otherCi: ['.gitlab-ci.yml', 'azure-pipelines.yml', '.circleci', 'Jenkinsfile', 'bitbucket-pipelines.yml'].filter(has),
  preCommit: ['.husky', 'lefthook.yml', '.lefthook.yml', '.pre-commit-config.yaml'].filter(has),
  docs: ['docs', '_docs', 'doc', 'README.md', 'ARCHITECTURE.md', 'CONTRIBUTING.md'].filter(has),
  gitattributes: has('.gitattributes'),
  envFiles: list('.').filter((n) => /^\.env(\..+)?$/i.test(n)),
};

// --- mode ------------------------------------------------------------------
let suggestedMode;
let reason;
if (sourceFiles === 0) {
  suggestedMode = 'NEW';
  reason = 'no source files found';
} else if (aiArtifacts.length > 0) {
  suggestedMode = 'EXISTING+AI';
  reason = `source files present and AI artifacts already exist (${aiArtifacts.join(', ')}): merge, never overwrite`;
} else {
  suggestedMode = 'EXISTING';
  reason = 'source files present, no AI configuration found';
}

const report = { target, exists: true, platform: process.platform, git: gitInfo, markers, languages: counts, sourceFiles, testFiles, scannedFiles: scanned, truncated: scanned > LIMIT, package: pkg, aiArtifacts, ai, graphCachePresent, harness, suggestedMode, reason };

if (asJson) {
  process.stdout.write(`${JSON.stringify(report, null, 2)}\n`);
} else {
  console.log(`target: ${target}`);
  console.log(`git: ${gitInfo.isRepo ? `${gitInfo.branch} (default: ${gitInfo.defaultBranch || 'unknown'}), ${gitInfo.dirtyFiles} dirty, ${gitInfo.commits} commits, .claude ignored: ${gitInfo.claudeDirIgnored}` : 'NOT a git repo'}`);
  console.log(`markers: ${markers.join(', ') || 'none'}`);
  console.log(`languages: ${Object.entries(counts).map(([k, v]) => `${k} ${v}`).join(', ') || 'none'} (tests: ${testFiles})`);
  if (pkg?.scripts) {
    const shown = pkg.scripts.slice(0, 30).join(', ');
    console.log(`scripts (${pkg.scripts.length}): ${shown}${pkg.scripts.length > 30 ? ', … (use --json for all)' : ''}`);
  }
  console.log(`AI artifacts: ${aiArtifacts.join(', ') || 'none'}${graphCachePresent ? ' · code-review-graph cache present' : ''}`);
  console.log(`CI: ${[...harness.ciWorkflows, ...harness.otherCi].join(', ') || 'none'} · pre-commit: ${harness.preCommit.join(', ') || 'none'} · env files: ${harness.envFiles.join(', ') || 'none'}`);
  console.log(`suggested mode: ${suggestedMode} (${reason})`);
}
