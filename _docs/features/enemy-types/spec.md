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
| **Movement** (how it moves and collides) | `ground` (gravity, walks on tiles), `flying` (no gravity, moves in 2D, keeps an altitude band), `static` (doesn't move) | `water` (stays inside water tiles) |
| **Attack** (how it hurts the player) | `melee` (today's attack box), `ranged` (fires a projectile), `dive` (flier swoops or smashes down, then recovers) | `contact`, multi-phase |

Optional per-species traits: **guard** (blocks some player hits, D8), a second melee clip
(`attack_2`) picked at random, a starting **sleep** state that wakes on proximity, and
colour **variants** that share one blueprint (Skeleton White/Yellow, Zombie 1/2, German
shepherd sizes).

## Standard clips

Every enemy is built from the same clip names, so the AI never asks which species it runs.
Required: `idle`, `walk`, `attack_1`, `hurt`, `death`. Optional, with a fallback when missing:
`run` → `alert` → `walk`, `attack_2` → `attack_1`, `guard` → none (the trait is off). Patrol
plays `walk` and chasing plays `run`. Names, fallbacks and the state → clip table live in
`internal/core/enemy_clips.go`.

## Two sprite formats, one game format (D3, D7)

- **Strip packs**: one PNG per clip, used as shipped (renamed when names collide, D4).
- **RPG Maker MV side-view battlers** (raw sheets in `art-src/battlers/`, not embedded): 9×6 grid, 18 motions × 3 frames
  (layout: `art-src/battlers/imagePattern.png`). An offline tool bakes each sheet into standard
  strips, chaining motions into longer clips:

| Clip | Battler motions chained | Frames |
|---|---|---|
| idle | Ready/Physical, ping-pong 1-2-3-2 | 4 |
| walk | Idle (MV's step) + bob at draw time | 4 |
| alert | Guard → Ready/Physical | 6 |
| run | walk frames, faster, + lean and bigger bob at draw time | 4 |
| attack_1 | Ready/Physical → Stabbing | 6 |
| attack_2 | Use Physical Skill → Swinging | 6 |
| guard | Guard (hold the middle frame while blocking) | 3 |
| hurt | Damage | 3 |
| death | Damage → Crisis/Danger → Dead (hold last) | 9 |
| shoot (ranged only) | Ready/Magical → Shooting | 6 |

Each sheet gets its own mapping, checked by eye: layouts vary (the goblin's "Escape" is a
tumble, not a run). The output faces right, has whole-pixel cells, and is scaled to draw size.

## Roster

Frame widths of strips are PNG width ÷ frame width (ASSUMED until checked by eye in the
enemy's task). Battler movement/attack choices are proposals (ASSUMED).

| Enemy | Source | Format | Move | Attack | Notes |
|---|---|---|---|---|---|
| Succubus (exists) | `succubus.png` | 270×360 grid | ground | melee | already wired |
| Skeleton White / Yellow | `Skeleton_White/`, `Skeleton_Yellow/` | strips 96×64 | ground | melee ×2 | idle 8, walk 10, attack1 10, attack2 9, hurt 5, die 13 (checked) |
| White Golem | `White_Golem/` | strips 90×64 | ground | melee, slow | idle 8, walk 10, attack 11, hurt 4, die 13; rename files (D4) |
| Mushroom | `Mushroom/` | strips 80×64 | ground | melee | idle 7, run 8, attack 10, hit 5, die 15 (checked); art faces left; stun clips later |
| Flying Demon | `flying_Demon/` | strips 81×71 | flying | melee | idle 4, flying 4, attack 8, hurt 4, death 7 |
| Smasher (Enemy3) | `Enemy3-Movement-In-Animation/` | strips 64×64 | flying | dive | idle 8, fly 8, smashStart 12, smashLoop 3, smashEnd 8, hit 4, die 17 |
| Goblin | `goblin sheet.png` | battler | ground | melee ×2, guard | |
| Kobold | `kobold_0000_red.png` | battler | ground | melee, guard | |
| Gnoll | `gnoll sheet.png` | battler | ground | melee ×2, guard | |
| Skeleton (battler) | `skelleton sheet.png` | battler | ground | melee ×2, guard | |
| Zombie 1 / 2 | `Zombies/` | battler | ground, slow | melee | variants |
| Troll | `troll_0000_green.png` | battler | ground, slow | melee ×2, guard | big |
| Werewolf | `Werewolf_0004_brown.png` | battler | ground, fast | melee ×2 | |
| Wolf | `wolf_0001_brown.png` | battler | ground, fast | melee | |
| German shepherd | `German shepard bundle/` | battler | ground, fast | melee | two sizes = variants |
| Rat | `Rat_0004_dark.png` | battler | ground | melee | faces left in the source |
| Crab | `crab sheet.png` | battler | ground | melee, guard | water later |
| Sahuagin | `sahuagin sheet.png` | battler | ground | melee | first `water` candidate |
| Slime (water) | `slime waterB sheet.png` | battler | ground | melee | |
| Turnip | `turnip sheet.png` | battler | ground | melee | |
| Bat | `Bat_0000_dark.png` | battler | flying | melee | |
| Training dummy | `training dummy sheet.png` | battler | static | none | practice target |
| *Boss: Striker* | `Striker/` | strips 128×96 | — | — | **out of scope** |
| *Boss: Orange Golem* | `Orange_Golem/` | strips 90×64 | — | — | **out of scope** |

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
   hits enemies on the ground. A guarding enemy takes no damage from a hit in front of it
   and plays `guard`; hits from behind land.
7. F3 shows body and attack boxes for every type. Native and WASM builds both pass.
8. Adding a species needs no edits to `enemy_ai.go`, `enemy_anim.go`, `level_enemies.go` or
   `cmd/main.go` unless it needs a new movement or attack part.
9. Baked battler strips are scaled to draw size: the enemy assets add well under the 15 MB
   the raw sheets (15 MB, now in `art-src/battlers/`) would add to the WASM download.

## Out of scope

Bosses (Striker, Orange Golem), mushroom stun, water enemies (the movement slot is designed
for them), newly drawn frames (only re-timing and draw-time motion), enemy sound,
difficulty scaling.
