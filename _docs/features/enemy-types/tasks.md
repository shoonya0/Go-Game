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
- [ ] 3. Standard clip names + fallbacks (`alert`/`run` → `walk`, `attack_2` → `attack_1`);
      chase plays `run` → check: `TestEnemyMissingClipFallsBack`
- [ ] 4. Per-frame `Hold` timing + draw-time motion (bob, lean, hit shake) → check:
      `TestClipHoldTiming`, manual look at the succubus
- [ ] 5. Extract `Mover` + `groundMover` from `enemy_ai.go`; succubus uses it → check:
      enemy tests unchanged and green
- [ ] 6. Extract `Attacker` + `meleeAttack` (damage on a `HitFrame`, `Box` for F3, random
      `attack_1`/`attack_2`) → check: `TestEnemyDetectsChasesAndAttacks` green, new
      `TestMeleeHitLandsOnHitFrame`
- [ ] 7. Registry validation test: unique markers, required clips present, speeds > 0,
      Mover/Attacker set → check: `TestBlueprintsValid`
- [ ] 8. `node scripts/verify.mjs --full`, `code-reviewer` agent, update docs (state.md,
      design.md trimmed to what the code doesn't show)

## Increment 2: battler baking tool

- [ ] 9. `tools/battler2strip` (stdlib Go): sheet + mapping JSON → right-facing strips,
      scaled to draw height, plus a printed blueprint snippet → check: `go test ./tools/...`
      (a cell-rect test with a non-integer cell height, e.g. 1298/6), depguard/layers still pass
- [ ] 10. Mapping for the goblin; bake; check every clip by eye (contact-sheet output) →
      check: baked strips look right, size printed
- [ ] 11. Decide what stays in `assets/`: keep raw sheets out of the `//go:embed` (move to
      `art-src/` or similar, ask the user) → check: WASM size before/after in `verify.mjs --full`

## Increment 3: ground roster (melee) + guard

- [ ] 12. Ledge probe in `groundMover`: patrol turns at platform edges → check:
      `TestGroundEnemyTurnsAtLedge`
- [ ] 13. Two-character level markers (D6) → check: `TestQuoteBareMarkers` extended, level loads
- [ ] 14. Guard trait (D8): `EnemyGuard`, `GuardChance`, front-only block → check:
      `TestGuardBlocksFrontHit`, `TestGuardIgnoresBackHit`
- [ ] 15. Goblin (first baked enemy, melee ×2 + guard) → check: manual look, F3 boxes
- [ ] 16. Rename `White_Golem/Golem_1_*` → `WhiteGolem_*` (D4), add White Golem → check:
      new `TestAssetBaseNamesUnique`
- [ ] 17. Skeleton White + Yellow strips as one blueprint with colour variants → check:
      `TestSkeletonVariantsShareBlueprint`, manual look
- [ ] 18. Mushroom (no stun yet) → check: manual look
- [ ] 19. Bake + add the other ground battlers, a few per session: kobold, gnoll, skeleton,
      zombies, troll, werewolf, wolf, German shepherd, rat, crab, sahuagin, slime, turnip →
      check: `TestBlueprintsValid`, manual look at each
- [ ] 20. Training dummy (`static` mover, no attack) → check: `TestStaticEnemyNeverMoves`
- [ ] 21. Place one of each in `level_1.json` (or a test level), web smoke check,
      `verify.mjs --full`, docs update

## Increment 4: flying movement

- [ ] 22. `flyingMover`: no gravity, hover height, leash, 2D chase, return to spawn,
      tile collision → check: `TestFlyingEnemyIgnoresGravity`, `TestFlyingEnemyLeash`
- [ ] 23. Ground wave skips enemies that aren't grounded (`Mover.Grounded`) → check:
      `TestGroundWaveMissesFlyers`
- [ ] 24. Flying Demon (strips) and Bat (battler) → check: manual look, F3 boxes

## Increment 5: more attack styles

- [ ] 25. `EnemyEffects` + `Projectile` (owned by `Game`, like ground waves); a ranged enemy
      (kobold or gnoll `shoot` clip) → check: `TestProjectileHitsPlayer`, `TestProjectileStopsAtTile`
- [ ] 26. `diveAttack` (start → loop → end); Smasher (Enemy3) → check: `TestDiveAttackPhases`, manual look
- [ ] 27. `verify.mjs --full`, web smoke check, `/code-review`, docs update, PR (ask before merge)

## Later (not planned yet)

- [ ] Bosses: Striker (`Striker/`, 128×96, dash/jump), Orange Golem (rename files first, D4)
- [ ] Mushroom stun attack (`Mushroom-AttackWithStun`, `Mushroom-Stun`): needs a player stun state
- [ ] `waterMover` (Sahuagin, crab, water slime are candidates)
- [ ] `readAsset` by full path instead of base name (project known issue #3)
