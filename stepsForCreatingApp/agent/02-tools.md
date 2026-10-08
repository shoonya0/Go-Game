# Phase 2: Tools (T1–T8)

Goal: every tool the chosen setup needs is **present, working and verified**, or
explicitly skipped by the user with a recorded consequence. **You never install
silently.** Each missing item goes through the Stop protocol (00-START §0.3), one item
at a time, required tools first.

Exact commands per OS: [`reference/install-commands.md`](./reference/install-commands.md).

---

## T1. Detect

```bash
node "<PLAYBOOK>/agent/scripts/check-tools.mjs"
```

Read `missing required` and `missing recommended`. Decide which items **this
project** needs:

| Tool | Needed when | Tier |
|---|---|---|
| git | always | required |
| node ≥ 20 + npm | always (hooks, scripts) | required |
| claude CLI | always (plugin/MCP checks) | required, already present if you're running |
| gh + `gh auth status` | Hosting = GitHub | required for GitHub; skip otherwise |
| uv **or** pipx **or** python+pip | to install code-review-graph | recommended (one installer is enough) |
| code-review-graph | always recommended (Depth ≠ minimal still includes it) | recommended |
| agent-browser + its Chrome | Platforms include Web or a desktop app with a web renderer (Electron/Tauri) | recommended for UI projects |
| language toolchain (python/go/cargo/dotnet/java …) | EXISTING projects in that language; NEW projects after the stack decision | required for that stack |

Write the needed list into the setup log (T1 row).

## T2. Required tools

For each missing required tool (git, node, gh), run the Stop protocol.

- **git (Windows):** Git for Windows also provides Git Bash, which Claude Code's Bash
  tool and the shell-form hooks use. After installing, restart Claude Code.
- **node:** after installing, restart Claude Code (PATH), then re-run T1.
- **gh:** install, then the **user** runs `! gh auth login` (interactive browser login,
  and a credential). Verify with `gh auth status`.

## T3. Python-tool installer (only if code-review-graph is missing)

Ask which installer to use (header `Installer`):
"uv (Recommended)" · "pipx" · "pip --user (python already present)".
Install the chosen one with the Stop protocol if it's missing.

## T4. code-review-graph

1. Stop protocol: install the CLI (`uv tool install code-review-graph`, or the pipx/pip
   alternative).
2. Verify: `code-review-graph --help` (and `uv tool update-shell` + restart if it's not
   on PATH).
3. **Do not run `code-review-graph install` or `build` yet.** Those run in phase 4 (H13),
   after the user approves the config changes they make.

## T5. agent-browser (UI projects only)

1. Stop protocol: `npm install -g agent-browser` (a global npm package).
2. Stop protocol again: `agent-browser install` downloads Chrome for Testing, a large
   download. Ask separately.
3. Verify: `agent-browser --version`.

## T6. Stack toolchain (EXISTING projects now; NEW projects after N3)

From `inspect-project.mjs` markers, check that the language toolchain runs (`python
--version`, `go version`, `cargo --version`, `dotnet --version`, `java -version`, …).
If one is missing, use the Stop protocol and link the official installer page. These
are often GUI or admin installs, so the **user** runs them.

## T7. Claude Code health

```bash
claude --version
claude doctor
claude plugin marketplace list
```

`claude doctor` is read-only. If `claude-plugins-official` isn't in the marketplace
list, use the Stop protocol for `claude plugin marketplace add anthropics/claude-plugins-official`.

## T8. Re-check and record

Re-run `check-tools.mjs`. Update the setup log's **Tools** table: DONE with the version,
SKIPPED with the consequence, or BLOCKED. A SKIPPED item disables its dependent steps:

| Skipped | Steps you then skip |
|---|---|
| gh | H11 GitHub CI push checks, `/install-github-app`, PR steps in the handoff |
| code-review-graph | E3/E4 graph tour (use Glob/Grep with a narrow scope instead), H13, the graph-update hook |
| agent-browser | H14 UI access (note: "UI verification is manual") |
| node | most of phase 4. Stop and tell the user the harness can't be completed without Node. |

✅ **Phase 2 checkpoint:** show the tools table and ask to continue. Then go to
`03-new-project.md` (NEW) or `04-existing-project.md` (EXISTING / EXISTING+AI).
