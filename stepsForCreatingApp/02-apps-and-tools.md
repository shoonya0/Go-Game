# 02: Apps and tools to install

Install in tiers. **Tier 0 and Tier 1 are enough to start.** Add the rest when a
trigger in the "Add when" column appears. Every MCP server, plugin and skill costs
some context, so install only what you'll use. Check what's loaded with `/context`.

Commands were checked on 2026-10-06. Items marked *(verify)* came from third-party
pages or change often.

---

## Tier 0: Foundation (required)

| Tool | What it's for | Install (Windows / macOS) |
|---|---|---|
| **Claude Code** (CLI; also Desktop app and VS Code/JetBrains extensions) | The coding agent | Windows PowerShell: `irm https://claude.ai/install.ps1 \| iex` · Windows CMD: `curl -fsSL https://claude.ai/install.cmd -o install.cmd && install.cmd && del install.cmd` · WinGet: `winget install Anthropic.ClaudeCode` · macOS: `curl -fsSL https://claude.ai/install.sh \| bash` or `brew install --cask claude-code`. Then run `claude --version` and `claude doctor`. |
| **A Claude plan** | Claude Code needs Pro, Max, Team, Enterprise or a Console (API) account. The free plan doesn't include it. | claude.ai, or platform.claude.com for API keys |
| **Git** (+ **Git for Windows** on Windows) | Version control is your undo button. On Windows, Git Bash powers Claude Code's Bash tool. Without it, Claude Code falls back to PowerShell. | git-scm.com · macOS: `xcode-select --install` or `brew install git` |
| **GitHub CLI `gh`** | PRs, issues, CI logs and review comments. Uses less context than an MCP server, and Claude already knows how to use it. | `winget install GitHub.cli` · `brew install gh` · then `gh auth login` |
| **Language runtime(s)** | Node.js LTS (also runs the hook scripts in `templates/`), plus Python with **uv** for Python-based tools | nodejs.org · `winget install astral-sh.uv` / `brew install uv` |
| **An editor for reviewing diffs** | You *read* code more than you write it | VS Code (+ Claude Code extension) or any IDE |

**Windows note.** On this machine, the existing global hooks in
`~/.claude/settings.json` use bash syntax and run fine, through Git Bash. For hooks you
share across macOS and Windows, prefer **Node scripts** (`node .claude/hooks/x.mjs`)
over shell one-liners. See [`templates/hooks/`](./templates/hooks/).

---

## Tier 1: Make the codebase readable for the agent (strongly recommended)

| Tool | Why | Install / use |
|---|---|---|
| **code-review-graph** (`tirth8205/code-review-graph`) | Local Tree-sitter graph of functions, classes, imports, calls and tests, stored in SQLite and served over MCP. It answers "who calls X?", "what's the blast radius?" and "which tests cover this?" without reading whole files. Natively uses it. | `pip install code-review-graph` (or `pipx install ...` / `uvx code-review-graph serve`), then `code-review-graph install` (sets up Claude Code and other agents) and `code-review-graph build`. Keep it fresh with `code-review-graph update` from a PostToolUse hook (see 03), and `code-review-graph detect-changes --brief` in pre-commit. Optional semantic search: `pip install "code-review-graph[embeddings]"`, then `code-review-graph embed`. Ships 4 skills: `explore-codebase`, `debug-issue`, `refactor-safely`, `review-changes`. |
| **LSP code-intelligence plugin** | Jump to definitions, find references, and see type errors right after each edit | `/plugin install typescript-lsp@claude-plugins-official` (also `pyright-lsp`, `gopls-lsp`, `rust-analyzer-lsp`, `jdtls-lsp`, `kotlin-lsp`, `csharp-lsp`, `swift-lsp`, `php-lsp`, `ruby-lsp`, `clangd-lsp`, `lua-lsp`). Some need the language server binary installed; see each plugin's README. |
| **Context7** (docs lookup) | Pulls **current, version-specific** library docs into context. It's the main fix for invented or outdated APIs. | `/plugin install context7@claude-plugins-official` (hosted MCP, no local Node needed) |
| **Serena** *(optional alternative)* | LSP-based semantic code navigation and refactoring over MCP | `/plugin install serena@claude-plugins-official` |

> Graph vs LSP: you don't strictly need both. The graph is best for **impact,
> callers, flows and test coverage across the repo**. LSP is best for **precise
> symbol lookup and live type errors**. On large codebases, use both.

---

## Tier 2: Let the agent verify its own work (required once you have a UI)

| Tool | Why | Install / use |
|---|---|---|
| **Your project's own checks** | Type checker, linter, formatter, unit/integration test runner, e2e runner. Without these there's no loop to close. | Wrap them in one command: `npm run verify` (see 07) |
| **agent-browser** (`vercel-labs/agent-browser`) | Rust CLI over Chrome DevTools Protocol (CDP). It returns accessibility-tree snapshots with short refs (`@e1`), so it uses far fewer tokens than raw HTML. It also drives **Electron** apps through their CDP port. Natively uses it. | `npm install -g agent-browser`, then `agent-browser install` (downloads Chrome for Testing). Load its skills each session with `agent-browser skills get core` (and `... electron`). Core loop: `open <url>` → `snapshot -i` → `click @e2` / `fill @e3 "text"` → **re-snapshot** → `console` / `errors` / `network requests` / `screenshot out.png`. For Electron, use `--cdp <port>` rather than `--auto-connect`. |
| **Playwright** (test runner and MCP) | Repeatable e2e tests in CI (`@playwright/test`), plus a browser the agent can drive over MCP | `npm i -D @playwright/test` · MCP: `/plugin install playwright@claude-plugins-official` |
| **Chrome DevTools MCP** | Performance traces, network and console of a live Chrome. Steinberger's one exception to "CLIs over MCP". | `/plugin install chrome-devtools-mcp@claude-plugins-official` |
| **Claude in Chrome** / screenshots | Visual checks against a design. Paste screenshots straight into the prompt; more than half of Steinberger's prompts include one. | Built-in Chrome integration (see Claude Code docs, "chrome") · Win+Shift+S / Cmd+Shift+4 |

---

## Tier 3: Code quality and entropy control (add from week 2)

| Tool | Catches | Use |
|---|---|---|
| **Formatter** (Prettier / Biome / ruff format / gofmt) | Style noise in diffs | Run it from a PostToolUse hook on every edit |
| **Linter** (ESLint / Biome / ruff / golangci-lint) | Bug patterns and conventions | Part of `verify`; fail CI on errors |
| **Architecture rules**: `dependency-cruiser` or `eslint-plugin-boundaries` (JS/TS), `import-linter` (Python) | Layer violations, such as "UI must not import the DB layer" | Write rules whose **error message says how to fix it**, as OpenAI did. Start in warn mode, then make them errors. |
| **knip** | Dead files, unused exports, unused dependencies | `npx knip` weekly, or in CI |
| **jscpd** | Copy-pasted code | `npx jscpd src` weekly |
| **Framework doctors** (e.g. **react-doctor**) | Framework-specific regressions | Natively runs it in pre-commit (advisory) and on PRs |
| **Pre-commit runner** (husky / lefthook / pre-commit) | Runs fast gates before every commit | Keep it **fast**; slow gates belong in CI |

---

## Tier 4: Security (required before anything public)

| Tool | Use |
|---|---|
| `/security-review` (built in) | Reviews the pending changes on the current branch |
| `security-guidance` plugin | Pattern warnings on edits, an LLM review of the diff on Stop, and a commit reviewer (injection, XSS, SSRF, secrets, and more) |
| `claude-security` plugin | Deep vulnerability scan at a chosen effort tier, where each finding is challenged before it's reported |
| `semgrep` / `sonarqube` plugins | Rule-based static analysis inside the agent loop |
| Secret scanner (e.g. gitleaks) *(verify)* | Blocks committed keys. Add it to pre-commit and CI. |
| `npm audit` / `pip-audit` / Dependabot | Known-vulnerable dependencies |
| **You** | Check every new dependency yourself: does it exist, is it reputable, is it maintained? (See slopsquatting in 01 §12.) |

---

## Tier 5: Team, automation and integrations (add when needed)

| Need | Tool | Install |
|---|---|---|
| `@claude` in issues and PRs, automated PR review, scheduled jobs | **Claude Code GitHub Action** (`anthropics/claude-code-action@v1`) | In Claude Code: `/install-github-app` (needs `gh` and repo admin). Secrets: `ANTHROPIC_API_KEY`, or `CLAUDE_CODE_OAUTH_TOKEN` from `claude setup-token`. |
| Scripts, CI and fan-out | Headless mode | `claude -p "<prompt>" --output-format json`, plus `--allowedTools`, `--permission-mode dontAsk` |
| Recurring agents | `/schedule` (cloud routines), `/loop` (local repeat) | Built in |
| Issue tracker | `gh` CLI, or the `linear` / `atlassian` / `notion` plugins | `/plugin install <name>@claude-plugins-official` |
| Production errors and analytics | `sentry`, `datadog`, `posthog`, `grafana-mcp` plugins | Same |
| Design handoff | `figma` plugin, `frontend-design` skill | Same |
| Platform and database | `vercel`, `netlify-skills`, `cloudflare`, `railway`, `render`, `supabase`, `neon`, `prisma`, `firebase`, `mongodb`, `stripe`, `expo`, `aws-*`, `postman` | Install **only** the ones your stack uses |
| Code review bots | `coderabbit`, `greptile` plugins, or Claude Code Review | Same |

---

## Tier 6: Optional productivity

| Tool | Why |
|---|---|
| **Voice dictation** (OS built-in: Win+H / macOS Dictation; Steinberger uses Wispr Flow) | Speaking gives longer, richer prompts with less effort |
| **Terminal panes** (Windows Terminal split panes / tmux / iTerm) | Watching 2 or 3 sessions at once when you go parallel |
| **A second model family** (e.g. Codex CLI) for spec or diff review | Different models catch different mistakes. Steinberger had GPT-5-Pro review his specs. Optional. |
| **Spec frameworks**: GitHub **Spec Kit**, **BMAD-METHOD**, OpenSpec, AWS **Kiro** IDE | A heavier, formal pipeline (constitution → spec → plan → tasks → implement). Spec Kit: `uv tool install specify-cli`, then `specify init <project> --integration <key>`. Its commands include `/speckit-constitution`, `-specify`, `-clarify`, `-plan`, `-tasks`, `-analyze`, `-implement` and `-converge`. *(verify the integration key for Claude Code in the Spec Kit integrations reference)* |

---

## Day-one install script (copy and adapt)

**Windows (PowerShell):**

```powershell
irm https://claude.ai/install.ps1 | iex          # Claude Code
winget install Git.Git GitHub.cli OpenJS.NodeJS.LTS astral-sh.uv
gh auth login
pip install code-review-graph                     # or: uv tool install code-review-graph
npm install -g agent-browser; agent-browser install
```

**macOS (zsh):**

```bash
curl -fsSL https://claude.ai/install.sh | bash    # or: brew install --cask claude-code
brew install git gh node uv
gh auth login
pip3 install code-review-graph                    # or: uv tool install code-review-graph
npm install -g agent-browser && agent-browser install
```

**Then, inside the project, in Claude Code:**

```text
/plugin install context7@claude-plugins-official
/plugin install typescript-lsp@claude-plugins-official      # your language's LSP
/plugin install superpowers@claude-plugins-official          # optional methodology, see 03
/reload-plugins
```

```bash
code-review-graph install && code-review-graph build
```

> `uv tool install code-review-graph` is the uv equivalent of `pipx install`; the
> README itself lists `pip`, `pipx` and `uvx`. *(verify)*
