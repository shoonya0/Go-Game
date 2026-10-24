# Enemy types: design

Increment 1 is in the code; this page keeps only what the code doesn't show (why it is
shaped this way, and the parts still to build). Reasons for decisions live in
[decisions.md](./decisions.md); clip names and fallbacks in [spec.md](./spec.md).

## Shape

A species is data: an `EnemyBlueprint` (`internal/core/enemy.go`, instances in
`enemy_types.go`) with clips, stats, a `Mover` and an `Attacker`. `UpdateEnemy`
(`enemy_ai.go`) is the shared brain: patrol → notice → chase → attack, hurt interrupts,
dead is final. `Attacker.InRange` decides "attack", so a ranged enemy attacks from far away
and a diver from above without a change to the brain.

## Built parts: what the code doesn't say

- **Mover** (`enemy_move.go`) is split into `Steer` (intent → velocity, called as the brain
  decides) and `Move` (integrate + collide), not one `Move(e, intent, qt)`: the brain's tests
  check velocity right after `decideBehavior`, and a dead enemy just steers to `moveHold` and
  keeps falling.
- **Attacker** (`enemy_attack.go`): the attack_1/attack_2 pick is a swappable func var (like
  `deathKneels`), not a per-enemy seeded `rand.Rand`. A swing lands once, on its clip's
  `HitFrame`. `AttackCooldown` can outlast a whole attack clip, so `Ready` holds a new swing
  back (the enemy idles, facing the player) until the cooldown ends by the hit frame;
  otherwise the skeleton whiffed every other swing (`TestMeleeEverySwingLands`). The
  succubus keeps `HitFrame` 0. `EnemyEffects` joins `Update` with projectiles (task 25).
- **Clips** (`enemy_clips.go`, `enemy_motion.go`): draw-time motion is in screen pixels and
  moves only the sprite; the collision box stays put (D7).

## Parts still to build

| Part | Notes |
|---|---|
| ledge probe in `groundMover` | task 12 |
| `flyingMover` | no gravity, `HoverHeight`, `Leash` (tiles), steers in X and Y, bobbing optional |
| `waterMover` | later; like flying but clamped to water tiles |
| `rangedAttack` | spawns a `Projectile` on `FireFrame`; `KeepDistance` |
| `diveAttack` | phases start → loop (falling) → end; uses up to 3 clips |

Part-specific runtime state (dive phase, fire timer, wake flag) goes in one small typed
field on `EnemyRuntime` per part, like `swing swingState`, not in a `map[string]any`; the
zero value means "idle".

## Battler baking (offline, `tools/battler2strip`)

Built (task 9). Usage and the mapping fields are in `tools/battler2strip/main.go` and
`bake.go` (`Mapping`, `ClipMapping`). A mapping per enemy lives next to the tool, not in the embed:

```json
{"name": "Goblin", "sheet": "art-src/battlers/goblin sheet.png",
 "facesLeft": true, "drawHeight": 120,
 "clips": {"attack_1": {"frames": [[0,1,0], [1,0,0], [1,0,1], [1,0,2]], "fps": 12, "hold": [1,1,2], "hitFrame": 2}}}
```

Each frame is `[motion block 0-2, row 0-5, frame 0-2]` (layout `art-src/battlers/imagePattern.png`).
Clip timing goes in the mapping only so the printed snippet is complete, which keeps
re-baking a single command. Cell edges are rounded from the exact fractional grid, so a
sheet whose size doesn't divide by 9 or 6 still tiles; each cell is then box-filtered to
one fixed frame size (a cell 1 px larger than its neighbour scales by < 0.5 %).

## Projectiles

`EnemyEffects` holds live projectiles. `Game` owns it, `updateEnemies` advances it, and
`Draw` draws it after enemies. Same pattern as the player's ground waves (`ground_wave.go`).
This is the only `cmd/main.go` change the feature needs.

## New states

- `EnemyGuard` (D8): a one-shot entered from `TakeDamage` when the blow comes from the
  facing side and a `GuardChance` roll passes. It deals no damage and hands back like hurt.
  The roll goes through a seeded `rand.Rand` held by the enemy, so tests are deterministic.
- `EnemyRun` is not a state: chase plays `run` and patrol plays `walk`.
- `EnemySleep` / `EnemyWake` only if a sleeping enemy is added (none in the current roster).
