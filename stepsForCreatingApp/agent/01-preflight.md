# Phase 1: Preflight (P1–P6)

Goal: know exactly **where**, **on what OS**, and **which mode** before touching
anything. This phase is read-only except for creating the setup log (P6).

---

## P1. Confirm the target folder

- Default `TARGET` = the session's working directory.
- If the working directory **is the playbook folder itself**, or a parent containing
  several projects, ask which folder is the target.
- Ask with AskUserQuestion (header `Target`). Options: "This folder: <cwd>",
  "A new folder (I'll give the path)", "Another existing folder (I'll give the path)".
- If TARGET is outside the session's working directory, ask the user to run
  `/add-dir <TARGET>`, or better, to restart Claude Code inside TARGET. Hooks and
  `.claude/` settings load from the folder Claude Code was started in.

✅ Check: `TARGET` is an absolute path and the user has confirmed it.

## P2. Identify the platform and Node

```bash
node -p "process.platform + ' ' + process.arch + ' node ' + process.version"
git --version
```

- If **Node is missing or older than 20**, STOP (Stop protocol, `install-commands.md#node`).
  Node is required: hooks and helper scripts are Node. Until it's installed, you can
  only do P1, P3 (by hand) and the interview (N2).
- Note which command tool you're using (Bash or PowerShell) and write commands for it.

## P3. Inspect the target (read-only)

```bash
node "<PLAYBOOK>/agent/scripts/inspect-project.mjs" "<TARGET>"
node "<PLAYBOOK>/agent/scripts/inspect-project.mjs" "<TARGET>" --json   # if you need details
```

Note: git state (repo? branch? default branch? dirty files?), languages, package
scripts, CI, pre-commit, env files, existing AI artifacts, whether `.gitignore` hides
`.claude/`, and the **suggested mode**.

## P4. Decide the mode (ask the user)

| Suggested | Meaning | Path |
|---|---|---|
| `NEW` | No source files (empty or missing folder) | `03-new-project.md` |
| `EXISTING` | Source code, no AI configuration | `04-existing-project.md` |
| `EXISTING+AI` | Source code plus CLAUDE.md / AGENTS.md / `.claude/` / `.mcp.json` / Cursor or Copilot rules | `04-existing-project.md` in **merge mode** |

Ask with AskUserQuestion (header `Mode`), putting the suggested mode first and marked
"(Recommended)". The user's answer wins. Example: a folder with only a README and a
prototype script may still be treated as NEW.

## P5. Ask the scope questions (one AskUserQuestion call, up to 4 questions)

1. **Platforms** (header `Platforms`, multiSelect): Web · Windows desktop · macOS
   desktop · Mobile (iOS/Android). This drives the CI OS matrix and the UI-testing
   setup.
2. **Depth** (header `Depth`): "Full harness (Recommended)" · "Minimal: CLAUDE.md,
   hooks, verify, graph only" · "Full, but no CI / no pre-commit".
3. **Methodology plugin** (header `Method`): "None, just the playbook loop
   (Recommended for first use)" · "superpowers" · "feature-dev" · "Spec Kit". Pick
   **one** at most (playbook 03 §3A).
4. **Git hosting** (header `Hosting`): GitHub · GitLab · Other/none. GitHub enables
   `gh`, Actions CI and `/install-github-app`.

Record the answers. Skip steps that don't apply (e.g. no CI for "Other/none").

## P6. Create or resume the setup log

- If `TARGET/docs/ai-harness-setup.md` exists → **resume** (00-START §0.4).
- Otherwise copy `PLAYBOOK/agent/templates/setup-log.md` to `TARGET/docs/ai-harness-setup.md`
  and fill in the header: date, TARGET, PLAYBOOK, OS, shell, MODE, LANG and the P5
  answers. Creating `docs/` and this one file needs no approval, since it's the agreed
  record. If `TARGET/docs/` already holds unrelated user documentation, ask whether to
  use `docs/` or `docs/ai/` instead.
- For a NEW project where TARGET doesn't exist yet: create the folder only after the
  user confirms (P1), then write the log.

✅ **Phase 1 checkpoint:** show the user a 5-line summary (target, OS, mode, platforms,
depth) and ask "Continue to tool checks?" Then open `02-tools.md`.
