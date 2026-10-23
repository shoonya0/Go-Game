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
| **Movement** (how it moves and collides) | `ground` (gravity, walks on tiles), `flying` (no gravity, moves in 2D, keeps an altitude band) | `water` (stays inside water tiles), `static` |
| **Attack** (how it hurts the player) | `melee` (today's attack box), `ranged` (fires a projectile), `dive` (flier swoops or smashes down, then recovers) | `contact`, multi-phase |

Optional per-species traits: a starting **sleep** state that wakes on proximity (Bat), a
second melee clip picked at random (Skeleton, Bat), and colour **variants** that share one
blueprint (Skeleton White/Yellow).

## Roster

Frame counts are PNG width ÷ frame width. Frame widths are inferred from the sheet sizes
(ASSUMED until each sheet is checked by eye in its task). Strip = one PNG per clip, one row.

| Enemy | Folder | Move | Attack | Frame (px) | Clips (frames) |
|---|---|---|---|---|---|
| Succubus (exists) | `succubus.png` | ground | melee | 270×360 grid | already wired |
| Snake | `1 Snake/` | ground | melee | 48×48 strips | idle 4, walk 4, attack 6, hurt 2, death 4 |
| Hyena | `2 Hyena/` | ground | melee | 48×48 strips | idle 4, walk 6, attack 6, hurt 2, death 6 |
| Scorpio | `3 Scorpio/` | ground | melee | 48×48 strips | idle 4, walk 4, attack 4, hurt 2, death 4 |
| Mummy | `5 Mummy/` | ground | melee | 48×48 strips | idle 4, walk 6, attack 6, hurt 2, death 6 |
| Deceased | `6 Deceased/` | ground | ranged (`Ball.png` 8×6) | 48×48 strips | idle 4, walk 6, attack 4, hurt 2, death 6 |
| Skeleton (White, Yellow) | `Skeleton_With_VFX/` | ground | melee ×2 | 96×64 strips | idle 8, walk 10, attack1 10, attack2 9, hurt 5, die 13 |
| White Golem | `White_Golem/` | ground | melee (slow, heavy) | 90×64 strips | idle 8, walk 10, attack 11, hurt 4, die 13 |
| Mushroom | `Mushroom/` | ground | melee | 80×64 strips | idle 7, run 8, attack 10, hit 5, die 15 (stun clips: later) |
| Slime | `slime/` | ground | melee/contact | 32×25 grid, 8 cols × 3 rows | rows to be mapped by eye |
| Vulture | `4 Vulture/` | flying | dive | 48×48 strips | idle 4, walk(fly) 4, attack 4, hurt 2, death 4 |
| Flying Demon | `flying_Demon/` | flying | melee | 81×71 strips | idle 4, flying 4, attack 8, hurt 4, death 7 |
| Bat | `Bat with VFX/` | flying | melee ×2 | 64×64 strips | sleep 3, wakeUp 16, idleFly 9, run 8, attack1 8, attack2 11, hurt 5, die 12 |
| Smasher (Enemy3) | `Enemy3-Movement-In-Animation/` | flying | dive (smash) | 64×64 strips | idle 8, fly 8, smashStart 12, smashLoop 3, smashEnd 8, hit 4, die 17 |
| *Boss: Striker* | `Striker/` | — | — | 128×96 | **out of scope** |
| *Boss: Orange Golem* | `Orange_Golem/` | — | — | 90×64 | **out of scope** |

The `<Name>.png` 48×48 files in the desert folders are single-frame portraits, not clips.

## Acceptance criteria

1. The succubus plays exactly as before after the refactor: existing enemy tests pass unchanged.
2. Every roster enemy (not the bosses) has a blueprint, a unique level marker, and spawns
   from `level_1.json` (or a test level) with all its clips drawing the right frames.
3. Ground enemies fall, stand on tiles, turn at walls and **don't walk off platform edges**
   while patrolling.
4. Flying enemies ignore gravity, collide with solid tiles, chase in 2D within a leash
   around their spawn, and return when the player leaves.
5. Ranged enemies fire a projectile that hurts the player on contact and is removed on
   hitting a tile or leaving the screen; dive enemies attack downward and recover.
6. Player attacks (J/I/K/L/O and U's ground wave) damage every type; the ground wave only
   hits enemies on the ground.
7. F3 shows body and attack boxes for every type. Native and WASM builds both pass.
8. Adding a species needs no edits to `enemy_ai.go`, `enemy_anim.go`, `level_enemies.go` or
   `cmd/main.go` unless it needs a new movement or attack part.

## Out of scope

Bosses (Striker, Orange Golem), mushroom stun, water enemies (no sprites yet; the movement
slot is designed for them), enemy sound, difficulty scaling.
