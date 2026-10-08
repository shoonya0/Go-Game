Here are the software products I found that were built entirely or mostly by AI, are useful for developers, and got wide adoption or had a large effect. I left out the many small "vibe-coded" apps that never caught on.

The strongest examples

1. Claude Code (Anthropic)
- How much AI wrote it: Anthropic engineers say about 90% of Claude Code is written with Claude Code. Anthropic's lead engineer has said he personally writes no code by hand any more.
- Developer-friendly? Yes. It's a terminal and IDE coding agent that you talk to in plain English.
- Impact: Its popularity was compared to ChatGPT's launch, and it became the reference point for "AI writes the code."
- Sources: Pragmatic Engineer (https://newsletter.pragmaticengineer.com/p/how-claude-code-is-built), Fortune (https://fortune.com/2026/01/29/100-percent-of-code-at-anthropic-and-openai-is-now-ai-written-boris-cherny-roon/), Medium (https://medium.com/codex/claude-code-2-1-the-ai-tool-that-writes-90-of-its-own-code-16e084fd5be0)

2. OpenClaw (formerly Clawdbot / Moltbot), by Peter Steinberger
- How much AI wrote it: Mostly OpenAI's Codex agents. Steinberger made more than 6,600 commits in one month and says openly, "I ship code I don't read."
- Developer-friendly? Yes. It's an open-source, self-hosted personal AI agent written in TypeScript and Swift. It connects to WhatsApp, Telegram, Slack, Discord, iMessage and other chat apps, and has a marketplace of about 1,700 community-built skills.
- Impact: GitHub called it the fastest-growing project in its history, with about 388K stars by August 2026. Steinberger then joined OpenAI.
- Sources: Wikipedia (https://en.wikipedia.org/wiki/OpenClaw), Pragmatic Engineer interview (https://newsletter.pragmaticengineer.com/p/the-creator-of-clawd-i-ship-code), Lex Fridman transcript (https://lexfridman.com/peter-steinberger-transcript/)

3. Claude Cowork (Anthropic)
- How much AI wrote it: All of it, according to Anthropic. When asked how much Claude Code wrote, the head of Claude Code answered "All of it." A team of four shipped it in about 10 days.
- Developer-friendly? Partly. It's "Claude Code without the terminal," aimed at non-coders. It runs on the same Agent SDK that developers can use.
- Impact: It's the most-cited example of release cycles shrinking from months to days.
- Sources: TechCrunch (https://www.techcrunch.com/2026/01/12/anthropics-new-cowork-tool-offers-claude-code-without-the-code/), The Algorithmic Bridge (https://www.thealgorithmicbridge.com/p/claude-code-coded-claude-cowork)

4. OpenAI Codex
- How much AI wrote it: OpenAI says almost 80% of its code is now written with Codex, and 97.9% of its employees use it.
- Developer-friendly? Yes. It's a coding agent available as a CLI, an IDE extension and a cloud service.
- Impact: It's one of the two main coding agents developers use, alongside Claude Code.
- Sources: Outlook Business (https://www.outlookbusiness.com/amp/story/deeptech/openais-strategy-head-in-india-reveals-codex-writes-80-of-companys-code), Metaintro (https://www.metaintro.com/blog/openai-employees-codex-ai-coding-preview-2026)

5. Base44, by Maor Shlomo
- How much AI wrote it: One developer built it with heavy AI help, with no employees, co-founder or investors at the start. Sources describe it as built solo, not with an exact AI percentage.
- Developer-friendly? It's an app builder where you describe the app in plain language. It's easy for beginners and quick for developers who want a prototype.
- Impact: It reached 250K users and about $189K a month in profit, then sold to Wix for $80M only 6 months after it started. It became the standard "solo founder plus AI" success story.
- Sources: Calcalist (https://www.calcalistech.com/ctechnews/article/s1iflnlelx), Lenny's Newsletter (https://www.lennysnewsletter.com/p/the-base44-bootstrapped-startup-success-story-maor-shlomo)

Excluded on purpose

- Lovable: about $500M in annual revenue and roughly 8M users, so it's huge for building apps with AI. But I couldn't find a source saying Lovable itself was mostly written by AI.
- Cursor, Cline, Roo Code, Kilo Code, Continue: popular AI coding tools, but I found no evidence that they were mostly written by AI.

Takeaway: the clearest cases where AI wrote most of a widely used product are Claude Code, OpenClaw, Cowork and Codex. All four are AI coding or agent tools that their makers built mostly with their own tools. OpenClaw is the most striking open-source case, and Base44 is the clearest commercial win.

Note on these numbers: the percentages are what the companies and creators themselves say. No one has audited them. In every case, humans still review the code and set up the checks it has to pass.

--

Web Search("Boris Cherny Claude Code workflow tips CLAUDE.md plan mode verification")
Did 1 search in 8s
Allowed by auto mode classifier

● Web Search("Peter Steinberger agentic engineering workflow close the loop tests multiple agents")
Did 1 search in 7s
Allowed by auto mode classifier

● Web Search("OpenAI harness engineering Codex million lines no human written code")
Did 1 search in 6s
Allowed by auto mode classifier

● Fetch(https://pub.towardsai.net/openais-harness-engineering-experiment-zero-manually-written-code-100a24ad04cf)
Received 0 bytes (403 Forbidden)
Allowed by auto mode classifier

● Fetch(https://steipete.me/posts/just-talk-to-it)
Received 27.6KB (200 OK)
Allowed by auto mode classifier

● Fetch(https://daily.dev/posts/how-openai-built-a-1m-line-app-with-codex-are-engineers-done--mplwtpxij)
Received 5.9KB (200 OK)
Allowed by auto mode classifier

● Fetch(https://gist.github.com/joyrexus/e20ead11b3df4de46ab32b4a7269abe0)
Received 134.8KB (200 OK)
Allowed by auto mode classifier

● Fetch(https://openai.com/index/harness-engineering/)
Received 0 bytes (403 Forbidden)
Allowed by auto mode classifier

● Web Search("OpenAI harness engineering AGENTS.md table of contents docs custom linters architecture layers garbage collection agent legibility")
Did 1 search in 7s
Allowed by auto mode classifier

The teams behind Claude Code, OpenClaw and OpenAI's 1M-line Codex experiment all describe the same idea: make the AI prove its work instead of hoping it doesn't hallucinate. The techniques below are ranked by how much they reduce mistakes for the effort they cost, best first.

Tier 1: biggest payoff, least effort

1. Let the AI check its own work (close the loop)
- What: give the agent a way to verify every change itself: run the tests, compile, lint, run the CLI, call the endpoint, or compare screenshots. Nothing counts as done until it's proven.
- Who: Boris Cherny, who created Claude Code, calls this "probably the most important tip" and says it improves output quality 2–3×. Peter Steinberger, who built OpenClaw, says agents become "dramatically more reliable" this way.
- Why it stops hallucination: made-up APIs and imaginary fixes fail the build or the tests right away, and the agent then fixes its own mistake.

2. Keep a living rules file (CLAUDE.md / AGENTS.md)
- What: a file the agent reads at the start of every session, listing commands, conventions and past mistakes. Each time you correct the AI, ask it to add the lesson to the file.
- Who: Cherny treats it as "a living set of rules." Steinberger's file is about 800 lines. OpenAI kept theirs to about 100 lines.
- Why it works: the same mistake doesn't come back in later sessions.

3. Plan before writing any code
- What: for anything non-trivial, use plan mode to agree on a plan with the AI first, then let it implement. With a good plan it usually gets it right in one go.
- Who: Cherny starts about 80% of his sessions in plan mode. Sometimes he has a second Claude review the plan as if it were a staff engineer.
- Why it works: wrong assumptions get caught before any code exists, which is the cheapest time to fix them.

4. Write tests right after each feature, in the same session
- What: ask the AI to write tests while the feature is still in its context.
- Who: Steinberger says tests written this way are better and catch real bugs.

Tier 2: high payoff, a bit of setup

5. Short entry file that points to deeper docs
- What: the main AGENTS.md is a table of contents of about 100 lines. Specs, plans and references live in a structured docs/ folder that the agent reads only when it needs them.
- Who: OpenAI's rule: "if agents can't see it [in the repo], it doesn't exist."
- Why it works: the agent's context stays small and accurate instead of overloaded.

6. Enforce the architecture with linters and structural tests
- What: custom lint rules, for example "this layer may only import from the layer next to it," that run on every PR. The error messages explain how to fix the problem.
- Who: this is how OpenAI kept 1M lines of AI-written code consistent.
- Why it works: the agent can't quietly break the design, and the error message tells it what to do instead.

7. Turn repeated work into skills, slash commands or hooks
- What: if you do something more than once a day, make it a command. Use PostToolUse hooks to auto-format after every edit.
- Who: the Claude Code team.
- Why it works: repeated work runs the same way every time, without improvising.

8. Keep each agent's context clean
- What: hand side tasks to subagents (Cherny's approach). Prefer CLIs over MCP servers so tools don't use up context (Steinberger's approach).
- Why it works: a long, cluttered context is a major cause of drift and made-up details.

9. Small, atomic commits made by the agent
- What: each agent commits only the files it edited.
- Who: Steinberger.
- Why it works: a bad change is easy to find and revert.

Tier 3: for scale and keeping quality up over time

10. Run several agents in parallel
- What: run 3–10 agents on separate tasks. Cherny uses one git worktree per session; Steinberger runs them in a terminal grid in the same folder.
- Why it works: this mostly adds speed rather than accuracy, and it only works safely once Tier 1 is in place.

11. Spend about 20% of the time on refactoring
- What: use jscpd to find duplicated code, knip to find dead code, and ESLint, then have the agents clean up.
- Who: Steinberger.
- Why it works: AI-written code piles up duplication fast, and a messy codebase makes future AI output worse.

12. Scheduled cleanup agents
- What: agents that run weekly to find docs that no longer match the code and broken architecture patterns.
- Who: OpenAI.

13. Give the agent screenshots and access to real logs
- What: more than half of Steinberger's prompts include a screenshot. Cherny points Claude at logs and lets it debug end to end.
- Why it works: real evidence beats a description the AI has to guess from.

The common thread

In all three cases the humans stopped writing code and started building the environment the AI works in: rules, checks and feedback loops. The AI still makes things up, but those mistakes get caught automatically before they're merged.

Your Natively repo already has several of these: CLAUDE.md, a docs/ folder and the code-review-graph. The biggest gaps worth checking are #1 (automatic verification loops) and #6 (architecture lint rules).

Sources:
- Boris Cherny team tips (gist) (https://gist.github.com/joyrexus/e20ead11b3df4de46ab32b4a7269abe0)
- 13 Claude Code tips from Boris Cherny (https://madappgang.com/blog/claude-code-tips-from-its-creator-boris-cherny/)
- How the creator of Claude Code uses it (https://paddo.dev/blog/how-boris-uses-claude-code/)
- Steinberger, "Just Talk To It" (https://steipete.me/posts/just-talk-to-it)
- Pragmatic Engineer: "I ship code I don't read" (https://newsletter.pragmaticengineer.com/p/the-creator-of-clawd-i-ship-code)
- Addy Osmani, Loop Engineering (https://addyosmani.com/blog/loop-engineering/)
- OpenAI harness engineering experiment (https://pub.towardsai.net/openais-harness-engineering-experiment-zero-manually-written-code-100a24ad04cf?gi=e09c5db29310)
- Martin Fowler, Harness engineering (https://martinfowler.com/articles/exploring-gen-ai/harness-engineering.html)
- harness-init scaffold (https://github.com/Gizele1/harness-init)
- Data Science Dojo, Harness engineering (https://datasciencedojo.com/blog/harness-engineering/)