---
name: setup-ai-harness
description: Set up the full AI-development harness (CLAUDE.md, hooks, skills, verify command, code graph, plugins, CI) in the current project, for a brand-new project or an existing codebase. Stops and asks before any install, login, overwrite, commit or push. Invoke manually with /setup-ai-harness.
disable-model-invocation: true
---

# Set up the AI harness

<!--
Install this launcher (optional):
  Windows:  copy this folder to  %USERPROFILE%\.claude\skills\setup-ai-harness\
  macOS:    copy this folder to  ~/.claude/skills/setup-ai-harness/
Then set PLAYBOOK below to the absolute path of your stepsForCreatingApp folder.
-->

PLAYBOOK = `D:\projects\native\natively-cluely-ai-assistant\stepsForCreatingApp`

1. If the PLAYBOOK path above doesn't exist on this machine, stop and ask the user for
   the correct path to the `stepsForCreatingApp` folder.
2. If reading files under PLAYBOOK is blocked because it's outside this project, ask the
   user to run `/add-dir <PLAYBOOK>`.
3. Read `<PLAYBOOK>/agent/00-START.md` and follow it exactly, with TARGET = the current
   project unless the user says otherwise. Extra instructions from the user: $ARGUMENTS
4. Hard rule from 00-START: **stop and ask before installing, authenticating,
   overwriting, committing or pushing anything.**
