# Stack matrix: what "verify" means per language

The agent uses this to fill `scripts/verify.mjs` STEPS or the package scripts (H2), the
format hook (H3), the linter (H8), the architecture rules (H9), the LSP plugin (H12)
and CI setup (H11). **Prefer tools the project already uses.** This table applies only
where nothing exists yet. Every tool added is a dev-dependency install, so use the
Stop protocol.

| Stack | Typecheck | Lint | Format (hook + check) | Unit test | E2E / UI | Architecture rules | LSP plugin | CI setup action |
|---|---|---|---|---|---|---|---|---|
| **TypeScript / JS** | `tsc --noEmit` | ESLint (flat config) or Biome | Prettier (`format-on-edit.mjs` as is) or Biome | Vitest / Jest / `node --test` | Playwright (`_electron` for Electron) | dependency-cruiser · eslint-plugin-boundaries | `typescript-lsp` | `actions/setup-node` |
| **Python** | pyright (or mypy) | `ruff check` | `ruff format` / `ruff format --check` | pytest | Playwright for Python | import-linter | `pyright-lsp` | `astral-sh/setup-uv` |
| **Go** | `go vet ./...` | golangci-lint | `gofmt -w` / `gofmt -l .` | `go test ./...` | Playwright (web) | depguard (golangci-lint) | `gopls-lsp` | `actions/setup-go` |
| **Rust** | `cargo check` | `cargo clippy -- -D warnings` | `cargo fmt` / `cargo fmt --check` | `cargo test` | n/a / Playwright (web) | module visibility (`pub(crate)`) + review | `rust-analyzer-lsp` | `dtolnay/rust-toolchain` *(verify)* |
| **Java** | compiler (`./gradlew compileJava` / `mvn compile`) | Checkstyle / SpotBugs / Error Prone | Spotless | JUnit (`./gradlew test`) | Playwright / Selenium | ArchUnit | `jdtls-lsp` | `actions/setup-java` |
| **Kotlin** | compiler | detekt | ktlint / Spotless | JUnit / Kotest | Espresso (Android) | ArchUnit / Konsist *(verify)* | `kotlin-lsp` | `actions/setup-java` |
| **C# / .NET** | `dotnet build -warnaserror` | .NET analyzers | `dotnet format` / `dotnet format --verify-no-changes` | `dotnet test` | Playwright for .NET | NetArchTest | `csharp-lsp` | `actions/setup-dotnet` |
| **Swift** | `swift build` / `xcodebuild` | SwiftLint | swift-format | XCTest / Swift Testing | XCUITest (manual-heavy) | module boundaries via SPM targets | `swift-lsp` | macOS runner + Xcode |
| **Ruby** | Sorbet (optional) | RuboCop | RuboCop `-a` | RSpec / Minitest | Capybara / Playwright | packwerk *(verify)* | `ruby-lsp` | `ruby/setup-ruby` |
| **PHP** | PHPStan / Psalm | PHPStan | PHP-CS-Fixer | PHPUnit / Pest | Playwright | Deptrac | `php-lsp` | `shivammathur/setup-php` |

## Notes

- **Desktop or multi-OS apps:** CI matrix `windows-latest` + `macos-latest` (+ Linux if
  shipped). Package scripts must work under cmd.exe: use `cross-env`, `rimraf` or Node
  scripts, never `VAR=x cmd` or `rm -rf`.
- **Electron:** native modules must be rebuilt for Electron's ABI. Tests that load them
  may need to run under Electron (`ELECTRON_RUN_AS_NODE=1 electron --test …`), as in
  Natively. Agent UI access goes through a dev-only CDP launcher (05-harness H14).
- **Monorepos:** run the checks per package (workspaces / turbo / nx). Put a CLAUDE.md in
  each package and use `.claude/rules/` with `paths:` for package rules.
- **AI or LLM features in the app:** add golden/invariant tests and a fixed-responder stub
  (playbook 07 §5). Never assert exact model text.
- **Formatter hook for non-JS stacks:** edit `format-on-edit.mjs` to call the stack's
  formatter on the single edited file, e.g. `spawnSync('ruff', ['format', file])`,
  `spawnSync('gofmt', ['-w', file])`, `spawnSync('rustfmt', [file])`. Keep it no-shell
  and exit 0 when the tool is missing.
