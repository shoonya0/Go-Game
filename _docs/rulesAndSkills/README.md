# Backend Engineering Skills

A set of [Claude Skills](https://docs.claude.com/en/docs/agents-and-tools/agent-skills/overview)
distilled from this Go game engine — a from-scratch 2D engine (quadtree, worker pool,
physics, WASM deploy).

Each skill captures a reusable backend practice with the real code that demonstrates it,
research-backed guidance, a checklist, and common pitfalls. They are written to show
concrete backend competency: concurrency, data structures, architecture and deployment.

| Skill | What it covers |
| ----- | -------------- |
| [`go-concurrent-worker-pool`](./go-concurrent-worker-pool/SKILL.md) | Bounded goroutine pools, channel fan-out/fan-in, graceful shutdown, CPU-aware sizing |
| [`go-spatial-partitioning`](./go-spatial-partitioning/SKILL.md) | Quadtree spatial index, O(log N) range queries, dynamic object→leaf updates |
| [`go-clean-architecture`](./go-clean-architecture/SKILL.md) | `cmd`/`internal` layout, consumer-defined interfaces, constructor injection |
| [`go-ebiten-wasm-deploy`](./go-ebiten-wasm-deploy/SKILL.md) | `go:embed` assets, WASM build, Go static server, multi-stage Docker, Render + CI |

## Using these skills

A skill is a directory containing a `SKILL.md` with YAML frontmatter (`name`, `description`)
and Markdown instructions. To make one available to Claude Code, copy its folder into a
skills directory:

```bash
# per-project
cp -r go-ebiten-wasm-deploy ../../.claude/skills/

# or per-user (all projects)
cp -r go-ebiten-wasm-deploy ~/.claude/skills/
```

Claude loads a skill's full instructions when a task matches its `description`, so keep
the description's trigger wording intact if you edit it.
