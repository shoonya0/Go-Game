# Enemy types: decisions

## D1. Compose behaviour from parts, not one Go type per species (2026-10-23)

- **Decision:** a species is a blueprint (data) plus a `Mover` and an `Attacker` value.
  The shared brain, animation and draw code never switch on the species.
- **Why:** ~14 enemies share almost everything. Per-species types would copy the state
  machine 14 times. One big `switch kind` would mean editing shared code for every enemy.
  The registry already works this way (`enemy_types.go`); this extends it to behaviour.
- **Rejected:** embedding a base struct per species (Go has no virtual dispatch, so the base
  can't call the overrides).

## D2. Two independent axes: movement and attack (2026-10-23)

- **Decision:** "flying / ground / water" is the movement axis only; attack style is a
  second axis.
- **Why:** the roster mixes them. The Vulture flies and dives, the Flying Demon flies and
  uses melee, the Deceased walks and shoots, the Bat flies and has two melee attacks. One
  axis would need a type for every combination.

## D3. Strips are the game's one sprite format; one sheet per clip (2026-10-23)

- **Decision:** a clip names its own sheet. Strip packs are used as shipped. Battler sheets
  are baked into strips offline (D7). Strips are **not** repacked into the battler format.
- **Why:** the battler format holds 3 frames per motion, and repacking strips into it would
  drop most frames (Golem die 13 → 3, Mushroom attack-with-stun 24 → 3). Strips also have
  clips (dash, wake-up, smash-loop) with no battler motion. Ebitengine batches draws well
  even across separate images.

## D4. Give every enemy asset a unique base name (2026-10-23)

- **Decision:** rename colliding files when an enemy is wired up, e.g.
  `White_Golem/Golem_1_*.png` → `WhiteGolem_*.png`, and `Orange_Golem` → `OrangeGolem_*`
  once the boss is added. **Done for the White Golem** (task 16); `TestAssetBaseNamesUnique`
  now fails on any duplicate name among the embedded files.
- **Why:** `readAsset` resolves by base name only (`internal/core/image.go`). Today both
  golem folders contain `Golem_1_*.png`, so one silently shadows the other. `Striker/` and
  `flying_Demon/` differ only by case (`idle.png` / `IDLE.png`). That works, since the embed
  FS is case-sensitive, but it is fragile, so prefix those too.
- **Rejected for now:** changing `readAsset` to full paths. It is the better fix, but it
  touches every asset call (project known issue #3). Worth its own refactor branch.

## D5. Bosses later (2026-10-23)

- **Decision:** Striker and Orange Golem are not built in this feature.
- **Why:** user request. Bosses need phases, arenas and health UI. The movement/attack
  parts should allow a boss to be a blueprint with a multi-phase `Attacker` later.

## D6. Level markers are a scarce resource (2026-10-23)

- **Note:** `isMarkerChar` allows about 20 punctuation characters, and with the battlers the
  roster has about 23 enemies plus variants, so they **will** run out. Do two-character markers
  (a change in `quoteBareMarkers`) before the battler roster lands (tasks.md), not letters,
  which clash with JSON literals.
- **Done (task 13):** a run of bare punctuation is one marker, up to `maxMarkerLen` (2)
  characters, so about 400 codes. Keep single characters for the existing enemies.

## D7. Bake battler sheets into strips with an offline tool (2026-10-23)

- **Decision:** `tools/battler2strip` (stdlib Go, not part of the game build) reads a
  battler sheet plus a mapping (clip → chain of motion cells, per-frame timing) and writes
  `<Enemy>_<clip>.png` strips facing right, with whole-pixel cells, scaled to draw size.
  Only the baked strips are embedded; the raw sheets stay out of `assets/` (or out of the
  embed).
- **Why:** one runtime format (D3). It fixes the 9 sheets whose height doesn't divide by 6
  and the mixed facing, and cuts the 15 MB of raw sheets (cells up to 435 px) to the drawn
  size, which matters for the WASM download (`Free Sprites/` alone is 15 MB).
- **Smoothness:** chaining motions gives 6–9 frame attacks and deaths. Draw-time motion
  (bob, lean, squash, hit shake) is applied in the game, not baked, so strip enemies get
  it too. New frames are not drawn: walk/run stay 3–4 frames.
- **Rejected:** reading battler sheets at runtime. Two formats in the game, raw sheets embedded.
- **Where the raw sheets live (task 11, 2026-10-24, user's choice):** `art-src/battlers/`,
  committed, outside `assets/` and the embed, so any machine can re-bake; `.dockerignore`
  keeps them out of the Render build context. Local WASM 51.2 → 36.2 MB.
- **Battler sheets deleted (2026-10-24, user's call):** none has a smooth walk or run
  (3 frames per motion, legs static), so the goblin slid. All 18 sheets removed from
  `art-src/battlers/`; restore one from git history if needed. Tasks 15, 19, 20 and the bat
  in 24 are dropped; `tools/battler2strip` stays until the enemy plan is redone.

## D8. Guard is an optional trait driven by a `guard` clip (2026-10-23)

- **Decision:** an enemy with a `guard` clip and `GuardChance > 0` may block a player hit
  that comes from the side it faces: no damage, `guard` plays, then it resumes. Hits from
  behind always land. Enemies without the clip never guard.
- **Why:** user request. Every battler sheet has a Guard motion; the strip packs don't, so it
  must be optional rather than part of the shared state machine's required clips.
- **Open (user to decide; no enemy guards yet):** an enemy mid-swing can't guard (my call),
  and one in its hurt reaction can guard the next combo hit (not decided). `blocks` in
  `enemy_guard.go` is where either rule lives.

## D9. Strip packs only; behaviour stays in parts (2026-10-24, user's call)

- **Decision:** every enemy comes from a strip pack in `assets/enemy/`. The battler sheets
  are deleted (see D7) and nothing is baked. An enemy is still a blueprint plus a `Mover`
  and an `Attacker` (D1, D2): the same clip-loading path and the same movement parts serve
  every species.
- **Why:** the goblin, the first baked battler, slid in its standing pose; the user turned
  it down. Strip packs all have real walk/run/fly frames, and one format means one loading
  path and no offline step.
- **Consequences:** D7's tool `tools/battler2strip` has nothing left to bake (keep or delete:
  user to decide). The battler-only plans are dropped: goblin, kobold/gnoll ranged enemy,
  training dummy (`static` mover), bat, water candidates. D8 guard works but no pack has a
  `guard` clip yet.
- **Succubus removed too (2026-10-24, user's call):** its animation didn't look good in the
  game. Its numbers and clips became the test-only `testEnemyBlueprint`, so the behaviour
  tests written against it still pin the same behaviour.
