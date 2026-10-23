# Enemy types: tasks

<!-- First unchecked box = where work resumes. One box per session-sized step, each
ending with its check. Tick only when the check passed. Don't renumber. -->

## Increment 1: modular foundation, no new enemies (design: [design.md](./design.md))

- [ ] 1. Characterisation tests for today's succubus first (patrol, chase, attack, hurt,
      death, ground landing) so the refactor can't change behaviour → check:
      `go test ./internal/core -run 'TestEnemy|TestBlueprint' -v`
- [ ] 2. Clips with their own sheet: `ClipDef` with sheet path + frame size + row; image
      cache keyed by path; succubus migrated → check: new `TestEnemyClipsFitTheirSheets`
      (decodes each sheet's config and checks frames × width fit), all enemy tests green
- [ ] 3. Optional clip fallback (`Alert → Patrol → Idle`) → check: `TestEnemyMissingClipFallsBack`
- [ ] 4. Extract `Mover` + `groundMover` from `enemy_ai.go`; succubus uses it → check:
      enemy tests unchanged and green
- [ ] 5. Extract `Attacker` + `meleeAttack` (damage on a `HitFrame`, `Box` for F3) → check:
      `TestEnemyDetectsChasesAndAttacks` green, new `TestMeleeHitLandsOnHitFrame`
- [ ] 6. Registry validation test: unique markers, every required clip present, speeds
      > 0, Mover/Attacker set → check: `TestBlueprintsValid`
- [ ] 7. `node scripts/verify.mjs --full`, `code-reviewer` agent, update docs (state.md,
      design.md trimmed to what the code doesn't show)

## Increment 2: ground roster (melee)

- [ ] 8. Ledge probe in `groundMover`: patrol turns at platform edges → check:
      `TestGroundEnemyTurnsAtLedge`
- [ ] 9. Desert four: Snake, Hyena, Scorpio, Mummy (48 px strips; pick scale against
      60 px tiles) + markers + embed check → check: `TestBlueprintsValid`,
      `TestAssetPathsAreEmbedded`, run `go run ./cmd` and look at each one
- [ ] 10. Rename `White_Golem/Golem_1_*` → `WhiteGolem_*` (D4), add White Golem → check:
      no duplicate base names under `assets/` (add `TestAssetBaseNamesUnique`)
- [ ] 11. Skeleton White + Yellow as one blueprint with colour variants; random pick
      between attack1/attack2 → check: `TestSkeletonVariantsShareBlueprint`, manual look
- [ ] 12. Mushroom (no stun yet) and Slime (map its 8×3 grid rows by eye) → check: manual look
- [ ] 13. Place one of each in `level_1.json` (or a test level), web smoke check,
      `verify.mjs --full`, docs update

## Increment 3: flying movement

- [ ] 14. `flyingMover`: no gravity, hover height, leash, 2D chase, return to spawn,
      tile collision → check: `TestFlyingEnemyIgnoresGravity`, `TestFlyingEnemyLeash`
- [ ] 15. Ground wave skips enemies that aren't grounded (`Mover.Grounded`) → check:
      `TestGroundWaveMissesFlyers`
- [ ] 16. Flying Demon (melee) → check: manual look, F3 boxes

## Increment 4: more attack styles

- [ ] 17. `EnemyEffects` + `Projectile` (owned by `Game`, like ground waves); Deceased with
      `rangedAttack` using `Ball.png` → check: `TestProjectileHitsPlayer`,
      `TestProjectileStopsAtTile`
- [ ] 18. `diveAttack` (start → loop → end); Vulture (swoop) → check: `TestDiveAttackPhases`
- [ ] 19. Smasher (Enemy3) using all three dive clips → check: manual look
- [ ] 20. `EnemySleep` / `EnemyWake` + Bat (sleep until the player is near, two attacks) →
      check: `TestBatWakesOnProximity`
- [ ] 21. `verify.mjs --full`, web smoke check, `/code-review`, docs update, PR (ask before merge)

## Later (not planned yet)

- [ ] Bosses: Striker (`Striker/`, 128×96, dash/jump), Orange Golem (rename files first, D4)
- [ ] Mushroom stun attack (`Mushroom-AttackWithStun`, `Mushroom-Stun`): needs a player stun state
- [ ] `waterMover` once there is a water-enemy sprite
- [ ] `readAsset` by full path instead of base name (project known issue #3)
- [ ] Two-character level markers if the ~20 single characters run out (D6)
