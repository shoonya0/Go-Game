# Enemy types: tasks

<!-- First unchecked box = where work resumes. One box per session-sized step, each
ending with its check. Tick only when the check passed. Don't renumber. -->

## Increment 1: modular foundation, no new enemies (design: [design.md](./design.md))

- [x] 1. Characterisation tests for today's succubus first (patrol, chase, attack, hurt,
      death, ground landing) so the refactor can't change behaviour → check:
      `go test ./internal/core -run 'TestEnemy|TestBlueprint' -v`
- [x] 2. Clips with their own sheet: `EnemyClip` with sheet path + frame size + row; image
      cache keyed by path; succubus migrated → check: new `TestEnemyClipsFitTheirSheets`
      (decodes each sheet's config and checks frames × width fit), all enemy tests green.
      It caught the succubus patrol clip running one column past its 8-column row; fixed
      to 7 frames (column 0 of that row is a front-facing stand, not a step)
- [x] 3. Standard clip names + fallbacks (`alert`/`run` → `walk`, `attack_2` → `attack_1`);
      chase plays `run` → check: `TestEnemyMissingClipFallsBack`
- [x] 4. Per-frame `Hold` timing + draw-time motion (bob, lean, hit shake) → check:
      `TestClipHoldTiming`, manual look at the succubus. The user dropped the succubus look
      (its art is weak); the look moved to the skeleton and mushroom (17, 18), pulled ahead
- [x] 5. Extract `Mover` + `groundMover` from `enemy_ai.go`; succubus uses it → check:
      enemy tests unchanged and green (19/19, no test file touched)
- [x] 6. Extract `Attacker` + `meleeAttack` (damage on a `HitFrame`, `Box` for F3, random
      `attack_1`/`attack_2`) → check: `TestEnemyDetectsChasesAndAttacks` green, new
      `TestMeleeHitLandsOnHitFrame` (mutation: ignoring HitFrame fails all 3 cases)
- [x] 7. Registry validation test: unique markers, required clips present, speeds > 0,
      Mover/Attacker set → check: `TestBlueprintsValid` (markers stay in
      `TestEnemyMarkersUniqueSingleChar`, now `TestEnemyMarkersUniqueAndBare`; mutation: no Attacker on the mushroom fails it)
- [x] 8. `node scripts/verify.mjs --full`, `code-reviewer` agent, update docs (state.md,
      design.md trimmed to what the code doesn't show). The reviewer found swings whiffing
      when the cooldown outlasts the clip; fixed test-first (`TestMeleeEverySwingLands`, `6c45e04`)

## Increment 2: battler baking tool

- [x] 9. `tools/battler2strip` (stdlib Go): sheet + mapping JSON → right-facing strips,
      scaled to draw height, plus a printed blueprint snippet → check: `go test ./tools/...`
      (a cell-rect test with a non-integer cell height, e.g. 1298/6), depguard/layers still pass.
      7 tests; depguard `tools` rule (stdlib only); mutations caught; test bake of the gnoll looked right
- [x] 10. Mapping for the goblin; bake; check every clip by eye (contact-sheet output) →
      check: baked strips look right, size printed
      `tools/battler2strip/mappings/goblin.json`; `-contact` added to the tool. 9 clips, 240×120
      frames, 423 KB; both attacks hit on frame 4. Strips in `assets/enemy/Goblin/` stay
      untracked until task 15 loads them. Hurt's flinch is subtle (hit shake covers it)
- [x] 11. Decide what stays in `assets/`: keep raw sheets out of the `//go:embed` (move to
      `art-src/` or similar, ask the user) → check: WASM size before/after in `verify.mjs --full`
      User chose `art-src/battlers/`, committed (D7). Local WASM 51.2 → 36.2 MB; re-bake from
      the new path is byte-identical; `art-src/` added to `.dockerignore`

## Increment 3: ground roster (melee) + guard

- [x] 12. Ledge probe in `groundMover`: patrol turns at platform edges → check:
      `TestGroundEnemyTurnsAtLedge`
- [x] 13. Two-character level markers (D6) → check: `TestQuoteBareMarkers` extended, level loads
- [x] 14. Guard trait (D8): `EnemyGuard`, `GuardChance`, front-only block → check:
      `TestGuardBlocksFrontHit`, `TestGuardIgnoresBackHit`
- [ ] ~~15. Goblin (first baked enemy, melee ×2 + guard) → check: manual look, F3 boxes.
      On hold: built (`849b30e`), then reverted (`12de28a`, `3b98b41`) because the sheet
      has no walk cycle and the user turned the sliding walk down~~ Dropped: battler sheets deleted (decisions)
- [x] 16. Rename `White_Golem/Golem_1_*` → `WhiteGolem_*` (D4), add White Golem → check:
      new `TestAssetBaseNamesUnique` (failed first on the 5 golem names). `%` on the platform at
      columns 44-49; `TestLevelEnemiesKeepTheirGround` added (`7dededa`..`d21c103`). User's native look pending
- [x] 17. Skeleton White + Yellow strips as one blueprint with colour variants → check:
      `TestSkeletonVariantsShareBlueprint`, manual look (user: "working"). Done early (`@`, `#` in level 1). `attack_2` plays once task 6 picks attacks
- [x] 18. Mushroom (no stun yet) → check: manual look (user: "working"). Done early (`&` in
      level 1, `FacesLeft`)
- [ ] ~~19. Bake + add the other ground battlers, a few per session: kobold, gnoll, skeleton,
      zombies, troll, werewolf, wolf, German shepherd, rat, crab, sahuagin, slime, turnip →
      check: `TestBlueprintsValid`, manual look at each~~ Dropped: battler sheets deleted
- [ ] ~~20. Training dummy (`static` mover, no attack) → check: `TestStaticEnemyNeverMoves`~~ Dropped: its sheet was a battler
- [x] 21a. Capability check (user request): `bp.Has(clip)` + a capability table → check:
      `TestBlueprintHasOwnArtOnly` (failed first, catches a `Has` that counts fallbacks), `TestEnemyCapabilities` (`1cbad37`, `9296da3`)
- [ ] 21. Place one of each in `level_1.json` (or a test level), web smoke check,
      `verify.mjs --full`, docs update

## Increment 4: flying movement

- [ ] 22. `flyingMover`: no gravity, hover height, leash, 2D chase, return to spawn,
      tile collision → check: `TestFlyingEnemyIgnoresGravity`, `TestFlyingEnemyLeash`
- [ ] 23. Ground wave skips enemies that aren't grounded (`Mover.Grounded`) → check:
      `TestGroundWaveMissesFlyers`
- [ ] 24. Flying Demon (strips) ~~and Bat (battler)~~ → check: manual look, F3 boxes

## Increment 5: more attack styles

- [ ] 25. `EnemyEffects` + `Projectile` (owned by `Game`, like ground waves); a ranged enemy
      (needs a strip pack with a shooting clip: none yet, D9) → check: `TestProjectileHitsPlayer`, `TestProjectileStopsAtTile`
- [ ] 26. `diveAttack` (start → loop → end); Smasher (Enemy3) → check: `TestDiveAttackPhases`, manual look
- [ ] 27. `verify.mjs --full`, web smoke check, `/code-review`, docs update, PR (ask before merge)

## Later (not planned yet)

- [ ] Bosses: Striker (`Striker/`, 96×96, attack 128×96, dash/jump), Orange Golem (rename files first, D4)
- [ ] Mushroom stun attack (`Mushroom-AttackWithStun`, `Mushroom-Stun`): needs a player stun state
- [ ] `waterMover` (no candidate since the battlers were dropped, D9)
- [ ] Delete `tools/battler2strip` or keep it (D9, user to decide)
- [ ] `readAsset` by full path instead of base name (project known issue #3)
