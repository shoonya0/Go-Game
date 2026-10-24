# Enemy types: spec

## Goal

Turn the sprite packs in `assets/enemy/` into enemies the player can fight. Behaviour is
composed from a few reusable parts, so that adding the next enemy (or a new kind of
movement such as swimming) is a blueprint entry plus, at most, one new part. The shared
update, draw and level code should not change.

## Classification

An enemy is described along two **independent** axes (see [decisions.md](./decisions.md) D2):

| Axis | Values now | Later |
|---|---|---|
| **Movement** (how it moves and collides) | `ground` (gravity, walks on tiles), `flying` (no gravity, moves in 2D, keeps an altitude band) | `static` (doesn't move), `water` (stays inside water tiles) |
| **Attack** (how it hurts the player) | `melee` (today's attack box), `ranged` (fires a projectile), `dive` (flier swoops or smashes down, then recovers) | `contact`, multi-phase |

Optional per-species traits: **guard** (blocks some player hits, D8; no current pack has a
`guard` clip, so no enemy guards yet), a second melee clip (`attack_2`) picked at random, a
starting **sleep** state that wakes on proximity, and colour **variants** that share one
blueprint (Skeleton White/Yellow).

## Standard clips

Every enemy is built from the same clip names, so the AI never asks which species it runs.
Required: `idle`, `walk`, `attack_1`, `hurt`, `death`. Optional, with a fallback when missing:
`run` → `alert` → `walk`, `attack_2` → `attack_1`, `guard` → none (the trait is off). Patrol
plays `walk` and chasing plays `run`. Names, fallbacks and the state → clip table live in
`internal/core/enemy_clips.go`.

The clips a species has are what it can do: shared code asks `bp.Has(ClipGuard)`, never
the species. There is no separate capability flag that could disagree with the art. To see
every enemy's capabilities: `go test ./internal/core -run TestEnemyCapabilities -v`.

## One sprite format: strip packs (D3, D9)

Every enemy is a strip pack: one PNG per clip, used as shipped (renamed when names collide,
D4). Each clip names its own sheet and frame size, so packs with different frame sizes
(Striker's attack is wider than its other clips) need no conversion. Every pack has real
multi-frame movement (walk/run/fly 4-10 frames), so no enemy needs faked walking.
The RPG Maker battler sheets were dropped on 2026-10-24 (D7, D9): 3 frames per motion and
no walk cycle.

## Roster

Frame sizes checked on 2026-10-24: every strip divides into whole frames, feet sit on a
steady row (fliers bob on purpose), and no frame is cut off.

| Enemy | Source | Format | Move | Attack | Notes |
|---|---|---|---|---|---|
| Skeleton White / Yellow | `Skeleton_White/`, `Skeleton_Yellow/` | strips 96×64 | ground | melee ×2 | idle 8, walk 10, attack1 10, attack2 9, hurt 5, die 13 (checked) |
| White Golem | `White_Golem/` | strips 90×64 | ground | melee, slow | idle 8, walk 10, attack 11 (slam on 6), hurt 4, die 13 (13th blank); files renamed `WhiteGolem_*` (D4); marker `%` |
| Mushroom | `Mushroom/` | strips 80×64 | ground | melee | idle 7, run 8, attack 10, hit 5, die 15 (checked); art faces left; stun clips later |
| Flying Demon | `flying_Demon/` | strips 81×71 | flying | melee | idle 4, flying 4, attack 8, hurt 4, death 7 |
| Smasher (Enemy3) | `Enemy3-Movement-In-Animation/` | strips 64×64 | flying | dive | idle 8, fly 8, smashStart 12, smashLoop 3, smashEnd 8, hit 4, die 17 |
| *Boss: Striker* | `Striker/` | strips 96×96 (attack 128×96) | — | — | **out of scope** |
| *Boss: Orange Golem* | `Orange_Golem/` | strips 90×64 | — | — | **out of scope** |

## Acceptance criteria

1. The succubus played exactly as before after the refactor: existing enemy tests passed unchanged. (Removed from the game on 2026-10-24; its numbers live on as the test-only `testEnemyBlueprint`.)
2. Every roster enemy (not the bosses) has a blueprint, a unique level marker, and spawns
   from `level_1.json` (or a test level) with all its clips drawing the right frames.
3. Ground enemies fall, stand on tiles, turn at walls and **don't walk off platform edges**
   while patrolling.
4. Flying enemies ignore gravity, collide with solid tiles, chase in 2D within a leash
   around their spawn, and return when the player leaves.
5. Ranged enemies fire a projectile that hurts the player on contact and is removed on
   hitting a tile or leaving the screen; dive enemies attack downward and recover.
6. Player attacks (J/I/K/L/O and U's ground wave) damage every type; the ground wave only
   hits enemies on the ground. A guarding enemy takes no damage from a hit in front of it
   and plays `guard`; hits from behind land.
7. F3 shows body and attack boxes for every type. Native and WASM builds both pass.
8. Adding a species needs no edits to `enemy_ai.go`, `enemy_anim.go`, `level_enemies.go` or
   `cmd/main.go` unless it needs a new movement or attack part.
9. ~~Baked battler strips are scaled to draw size~~ Dropped with the battler sheets (D9).

## Out of scope

Bosses (Striker, Orange Golem), mushroom stun, water enemies (the movement slot is designed
for them), newly drawn frames (only re-timing and draw-time motion), enemy sound,
difficulty scaling.
