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

## D3. Keep the sprite strips as shipped; one sheet per clip (2026-10-23)

- **Decision:** a clip names its own sheet. Don't repack strips into one grid the way the
  succubus sheet was repacked.
- **Why:** repacking is manual work for each enemy and has to be redone when an asset
  changes. Ebitengine batches draws well even across separate images, and the strips are small
  (the largest folder is under 70 KB).

## D4. Give every enemy asset a unique base name (2026-10-23)

- **Decision:** rename colliding files when an enemy is wired up, e.g.
  `White_Golem/Golem_1_*.png` → `WhiteGolem_*.png`, and `Orange_Golem` → `OrangeGolem_*`
  once the boss is added.
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

- **Note:** `isMarkerChar` allows about 20 punctuation characters, and the roster needs about 14
  plus variants. Assign them in [spec.md](./spec.md)'s order when each enemy lands. If they run
  out, move to two-character markers (a change in `quoteBareMarkers`), not letters, which
  clash with JSON literals.
