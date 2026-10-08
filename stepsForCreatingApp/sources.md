# Sources

Accessed 2026-10-06 unless noted. **Reliability key:** 🟢 official docs or the tool's
own README · 🟡 a practitioner's first-hand account · 🟠 a secondary summary or
self-reported claim, not independently audited.

## Official Claude Code documentation 🟢

- Best practices: https://code.claude.com/docs/en/best-practices (verification first;
  explore → plan → code → commit; CLAUDE.md include/exclude table; interview → SPEC.md;
  writer/reviewer; fan-out; failure patterns)
- Memory, CLAUDE.md and AGENTS.md: https://code.claude.com/docs/en/memory (AGENTS.md is
  read when no CLAUDE.md exists; `.claude/rules/` with `paths:`; auto memory)
- Extend Claude Code (features overview): https://code.claude.com/docs/en/features-overview
  (feature-to-goal table; "build your setup over time" triggers; context cost by feature;
  keep CLAUDE.md under 200 lines)
- Hooks reference: https://code.claude.com/docs/en/hooks (event list; exit code 2 /
  `decision: block`; `permissionDecision: deny`)
- GitHub Actions: https://code.claude.com/docs/en/github-actions (`/install-github-app`,
  `anthropics/claude-code-action@v1`, `@claude`, skills as prompts)
- Setup / install: https://code.claude.com/docs/en/setup (native installer, WinGet,
  Homebrew; Git for Windows provides Git Bash for the Bash tool)
- Docs index for anything else: https://code.claude.com/docs/llms.txt

## Practitioners

- Boris Cherny and the Claude Code team, tips (gist) 🟡:
  https://gist.github.com/joyrexus/e20ead11b3df4de46ab32b4a7269abe0 (3–5 parallel
  worktrees; plan mode plus a staff-engineer review; living CLAUDE.md; slash commands for
  daily tasks; PostToolUse formatting; `/permissions` instead of skipping permissions)
- Other Cherny write-ups 🟠: https://madappgang.com/blog/claude-code-tips-from-its-creator-boris-cherny/ ·
  https://paddo.dev/blog/how-boris-uses-claude-code/
- Peter Steinberger, "Just Talk To It" 🟡: https://steipete.me/posts/just-talk-to-it
  (3–8 agents in the same folder; screenshots in about 50% of prompts; an AGENTS.md of
  about 800 lines; CLIs over MCP; about 20% of time on refactoring with jscpd and knip;
  atomic commits; skeptical of subagents and big specs)
- Pragmatic Engineer, "I ship code I don't read" 🟡:
  https://newsletter.pragmaticengineer.com/p/the-creator-of-clawd-i-ship-code
- OpenAI, "Harness engineering: leveraging Codex in an agent-first world" 🟡 (the
  original returned 403 when fetched; read through these mirrors):
  https://www.engineering.fyi/article/harness-engineering-leveraging-codex-in-an-agent-first-world ·
  https://openai.com/index/harness-engineering/ (3→7 engineers, 5 months, about 1M lines,
  about 1,500 PRs; an AGENTS.md of about 100 lines as a map; layered architecture enforced
  by custom linters with remediation messages; app bootable per worktree; CDP for the
  agent; agent-queryable logs, metrics and traces; "garbage collection" agents) 🟠
  numbers are self-reported
- Martin Fowler, "Harness engineering" 🟡:
  https://martinfowler.com/articles/exploring-gen-ai/harness-engineering.html (guides vs
  sensors; computational vs inferential; caveats: the behavior-harness gap, harnessability
  of legacy code, the human role)
- Addy Osmani, "Loop engineering" 🟡: https://addyosmani.com/blog/loop-engineering/
  (automations, worktrees, skills, connectors, maker/checker; comprehension debt;
  cognitive surrender)

## Tools 🟢 (READMEs)

- code-review-graph: https://github.com/tirth8205/code-review-graph
- agent-browser: https://github.com/vercel-labs/agent-browser
- GitHub Spec Kit: https://github.com/github/spec-kit (integrations reference:
  https://github.github.io/spec-kit/reference/integrations.html)
- superpowers: https://github.com/obra/superpowers (also listed in the official
  marketplace as `superpowers`)
- claude-code-action: https://github.com/anthropics/claude-code-action
- Anthropic skills repository: https://github.com/anthropics/skills · Agent Skills overview:
  https://platform.claude.com/docs/en/agents-and-tools/agent-skills/overview
- Official plugin marketplace: https://github.com/anthropics/claude-plugins-official
  (plugin names verified in the local copy at
  `~/.claude/plugins/marketplaces/claude-plugins-official/.claude-plugin/marketplace.json`)
- levnikolaevich/claude-code-skills (vendored at `_docs/claude-code-skills/`)

## Spec-driven development overviews 🟠

- https://www.augmentcode.com/tools/best-spec-driven-development-tools
- https://reenbit.com/bmad-vs-spec-kit-vs-openspec-choosing-your-spec-driven-ai-framework/
- https://dev.to/krlz/spec-driven-development-in-2026-what-it-is-the-tooling-and-how-teams-actually-use-it-2fk2

## Security of AI-generated code

- Package hallucination (USENIX Security 2025; about 19.7% of recommended packages
  didn't exist across 576k samples from 16 LLMs) 🟢 peer-reviewed; summarized at
  https://nhimg.org/articles/slopsquatting-exposes-a-new-software-supply-chain-risk-for-ai-coding/ 🟠
- Veracode 2025 GenAI code security report (about 45% of samples had OWASP Top 10 flaws) 🟠
  as cited in https://dev.to/coppersundev/is-ai-generated-code-buggier-the-2025-26-data-5acf
- CSA research note on the rise in CVEs from AI-generated code 🟠:
  https://labs.cloudsecurityalliance.org/research/csa-research-note-ai-generated-code-vulnerability-surge-2026/

## Secondary guides consulted (for cross-checking only) 🟠

- https://github.com/shanraisshan/claude-code-best-practice
- https://alexop.dev/posts/understanding-claude-code-full-stack/
- https://www.developersdigest.tech/blog/claude-code-superpowers-plugin-guide

## This repository (first-hand inspection, commit `fcbfdd2`) 🟢

`CLAUDE.md`, `_docs/README.md`, `_docs/04-phased-roadmap.md`, `_docs/current-state.md`,
`_docs/phase-c6-status.md`, `_docs/test/README.md`, `docs/superpowers/plans/2026-09-23-fast-model-picker.md`,
`docs/retrieval-handoff/00-START-HERE.md`, `scripts/dev-agent.mjs`, `.husky/pre-commit`,
`.github/workflows/build-smoke.yml`, `.github/workflows/react-doctor.yml`, `.mcp.json`,
`package.json`, `~/.claude/settings.json`, and the auto-memory notes.

## Earlier research

- [`previous-research.md`](./previous-research.md): products built mostly by AI
  (Claude Code, OpenClaw, Cowork, Codex, Base44) and the ranked techniques list. Its
  percentages are self-reported 🟠.
