# Install and verify commands, per tool and per OS

The agent quotes these inside the Stop protocol (00-START §0.3). **Who runs it** says
whether the agent may offer "You run it", or whether the user must run it because it's
interactive, needs admin rights, or involves a login.

In Claude Code the user can type `! <command>` to run a command in this session, so its
output lands in the conversation. After any install that changes PATH, especially on
Windows, **restart Claude Code** (`claude --continue`) before verifying.

Checked 2026-10-06. *(verify)* marks commands taken from third-party pages.

---

## node

| OS | Install | Who runs |
|---|---|---|
| Windows | `winget install OpenJS.NodeJS.LTS` | User (installer, may show UAC) |
| macOS | `brew install node` (or nvm / fnm / the nodejs.org installer) | User |

Verify: `node --version` (≥ 20) and `npm --version`.

## git (+ Git Bash on Windows)

| OS | Install | Who runs |
|---|---|---|
| Windows | `winget install --id Git.Git -e` | User |
| macOS | `xcode-select --install` or `brew install git` | User |

Verify: `git --version`. On Windows, Claude Code uses Git Bash for its Bash tool when
Git for Windows is installed. If it can't find bash, set `CLAUDE_CODE_GIT_BASH_PATH` in
settings `env`.

## gh (GitHub CLI) + login

| OS | Install | Who runs |
|---|---|---|
| Windows | `winget install --id GitHub.cli -e` | User |
| macOS | `brew install gh` | User, or agent after OK |

Then the login, which is interactive and handles a credential: **the user** runs
`! gh auth login`. Verify: `gh --version` and `gh auth status` (exit 0 means logged in).

## uv (preferred installer for Python CLI tools)

| OS | Install | Who runs |
|---|---|---|
| Windows | `winget install --id astral-sh.uv -e` or `powershell -ExecutionPolicy ByPass -c "irm https://astral.sh/uv/install.ps1 \| iex"` | User |
| macOS | `brew install uv` or `curl -LsSf https://astral.sh/uv/install.sh \| sh` | User, or agent after OK (brew) |

Verify: `uv --version`. If tools installed with `uv tool` aren't on PATH:
`uv tool update-shell`, then restart the terminal.

## pipx (alternative)

Windows: `python -m pip install --user pipx` then `python -m pipx ensurepath` ·
macOS: `brew install pipx` then `pipx ensurepath`. Verify: `pipx --version`.

## code-review-graph

| Installer | Command | Who runs |
|---|---|---|
| uv | `uv tool install code-review-graph` | Agent after OK, or user |
| pipx | `pipx install code-review-graph` | Agent after OK, or user |
| pip | `pip install --user code-review-graph` (Windows: `py -m pip install --user code-review-graph`) | Agent after OK, or user |

Verify: `code-review-graph --help`.
Project integration happens later, in H13, and asks again:
`code-review-graph install --platform claude-code --repo <TARGET> --dry-run` (preview) → then for real
(recommended flags: `--no-instructions --no-hooks`) → `code-review-graph build --repo <TARGET>`.
Optional semantic search (large download): `uv tool install --force "code-review-graph[embeddings]"`
then `code-review-graph embed --repo <TARGET>`.

## agent-browser (UI projects)

| Step | Command | Who runs |
|---|---|---|
| CLI | `npm install -g agent-browser` (macOS also: `brew install agent-browser`) | Agent after OK, or user |
| Browser | `agent-browser install` (downloads Chrome for Testing, large) | Agent after a **separate** OK, or user |

Verify: `agent-browser --version`. Each session: `agent-browser skills get core` (and
`agent-browser skills get electron` for Electron). Never use `--auto-connect`; pass
`--cdp <port>`.

## Claude Code plugins

| Action | Command | Who runs |
|---|---|---|
| Check availability | `claude plugin list --available --json` | Agent (read-only) |
| Install (shared with the team) | `claude plugin install <name>@claude-plugins-official --scope project` | Agent after OK, or user (`/plugin install …` in the session) |
| Load without restarting | `/reload-plugins` | User (slash command) |
| Verify / cost | `claude plugin list --json`, `claude plugin details <name>@claude-plugins-official` | Agent |
| Missing official marketplace | `claude plugin marketplace add anthropics/claude-plugins-official` | Agent after OK |

LSP plugins may also need a language-server binary; read the plugin README. Installing
that binary is another Stop protocol.

## MCP servers

List or inspect: `claude mcp list`, `claude mcp get <name>` (agent, read-only).
Project servers from `.mcp.json` stay "Pending approval" until the user approves them
in a session (`/mcp` or the startup prompt).

## Project dependencies (always the Stop protocol; list each package and why)

| Stack | Example |
|---|---|
| npm | `npm install -D prettier eslint @eslint/js typescript-eslint dependency-cruiser husky` |
| pnpm / yarn / bun | same packages with `pnpm add -D` / `yarn add -D` / `bun add -d` (use the project's lockfile manager) |
| uv | `uv add --dev ruff pyright pytest import-linter pre-commit` |
| Go | `go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest` *(verify)* |
| Rust | `rustup component add clippy rustfmt` |

Before adding **any** new package name, confirm that it exists and is the real,
maintained package: check the registry page, downloads and repository. AI-suggested
package names can be made up ("slopsquatting", playbook 01 §12).

## Optional

| Tool | Install | Notes |
|---|---|---|
| gitleaks (secret scanner) | Windows: `winget install gitleaks` *(verify the id)* · macOS: `brew install gitleaks` | User runs |
| Claude GitHub app | `/install-github-app` in a session | User runs (interactive, repo admin) |
