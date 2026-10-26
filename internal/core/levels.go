package core

import (
	"fmt"
	"log"
	"slices"
	"strings"
)

// LevelDef says where a named level comes from: the JSON bundle and key its
// grid is read from, its background image and its size in world pixels.
type LevelDef struct {
	Title, Time   string // shown in the travel menu: "Deep Forest", "Day"
	File, Key     string
	Background    string
	Width, Height float64
	Tiles         []TileType    // the only tile types its grid may use (levels spec)
	KeepApart     [][2]TileType // tile types that must never touch, even diagonally
}

// StartLevel is the level a new game begins in: deepForest by day. Its name
// predates the day/night split and stays, since saves store it.
const StartLevel = "deepForest1"

// DeepForestNight is deepForest's night sub-level (levels decision D1).
const DeepForestNight = "deepForestNight"

// The other sub-levels, by level in the spec's order (levels decision D1).
const (
	AquaticAbyssDay     = "aquaticAbyssDay"
	AquaticAbyssNight   = "aquaticAbyssNight"
	BoulderBastionDay   = "boulderBastionDay"
	BoulderBastionNight = "boulderBastionNight"
	CargoYardChaosDay   = "cargoYardChaosDay"
	CargoYardChaosNight = "cargoYardChaosNight"
)

// subLevel is a sub-level laid out the usual way: the grid in
// assets/levels/<level>/<level>_<time>.json under the key time ("day",
// "night"), the background in its background folder with the same base name.
func subLevel(level, title, time, bgExt string, tiles ...TileType) LevelDef {
	dir := "../assets/levels/" + level + "/"
	base := level + "_" + time
	return LevelDef{
		Title: title, Time: strings.ToUpper(time[:1]) + time[1:],
		File: dir + base + ".json", Key: time,
		Background: dir + "background/" + base + bgExt,
		Width:      Level_1_Width, Height: Level_1_Height,
		Tiles: tiles,
	}
}

// boulderBastion keeps larva away from grass (levels spec).
func boulderBastion(time, bgExt string) LevelDef {
	def := subLevel("boulderBastion", "Boulder Bastion", time, bgExt, Rock, Grass, Larva)
	def.KeepApart = [][2]TileType{{Larva, Grass}}
	return def
}

// Levels lists every level by name. Travel between portals switches by name.
var Levels = map[string]LevelDef{
	StartLevel: {
		Title: "Deep Forest", Time: "Day",
		File: LevelFile, Key: DefaultLevelName, Background: Background_1,
		Width: Level_1_Width, Height: Level_1_Height,
		Tiles: []TileType{Grass, Water, Rock},
	},
	DeepForestNight: {
		Title: "Deep Forest", Time: "Night",
		File: "../assets/levels/deepForest/deepForest_night.json", Key: "night",
		Background: "../assets/levels/deepForest/background/deepForest_night.jpg",
		Width:      Level_1_Width, Height: Level_1_Height,
		Tiles: []TileType{Grass, Water, Rock},
	},
	AquaticAbyssDay: {
		Title: "Aquatic Abyss", Time: "Day",
		File: "../assets/levels/aquaticAbyss/aquaticAbyss_day.json", Key: "day",
		Background: "../assets/levels/aquaticAbyss/background/aquaticAbyss_day.png",
		Width:      Level_1_Width, Height: Level_1_Height,
		Tiles: []TileType{Water, Rock, Grass},
	},
	AquaticAbyssNight:   subLevel("aquaticAbyss", "Aquatic Abyss", "night", ".png", Water, Rock, Grass),
	BoulderBastionDay:   boulderBastion("day", ".jpg"),
	BoulderBastionNight: boulderBastion("night", ".jpg"),
	CargoYardChaosDay:   subLevel("cargoYardChaos", "Cargo Yard Chaos", "day", ".png", Wood, Grass, Water, Rock),
	CargoYardChaosNight: subLevel("cargoYardChaos", "Cargo Yard Chaos", "night", ".png", Wood, Grass, Water, Rock),
}

// LevelOrder is the order levels open in: saving the last unused portal of
// one takes the player to the next (levels decision D5).
var LevelOrder = []string{
	StartLevel, DeepForestNight,
	AquaticAbyssDay, AquaticAbyssNight,
	BoulderBastionDay, BoulderBastionNight,
	CargoYardChaosDay, CargoYardChaosNight,
}

// levelPortals caches each level's portal cells in reading order, read once
// for portalNumber.
var levelPortals = map[string][]Position{}

// portalNumber is ref's place among its level's portals in reading order,
// from 1 (1 is the entry portal), or 0 when ref's level isn't in Levels.
func portalNumber(ref PortalRef) int {
	cells, ok := levelPortals[ref.Level]
	if !ok {
		def, known := Levels[ref.Level]
		if !known {
			return 0
		}
		for _, c := range LoadLevel(def.File, def.Key).Portals {
			cells = append(cells, c.Pos)
		}
		levelPortals[ref.Level] = cells
	}
	return slices.Index(cells, ref.Cell) + 1
}

// World is the level the player is in: its platforms (indexed in Quadtree)
// and its enemies. It remembers the enemies of levels the player left.
type World struct {
	Name     string
	Def      LevelDef
	Level    Level
	Quadtree *DynamicQuadtree
	Enemies  []*EnemyRuntime

	left   map[string][]*EnemyRuntime // enemies of each level left, as they were (for its bosses)
	border []Platform                 // the current level's invisible wall (levelBorder)
}

// NewWorld loads the named level with fresh enemies. WorldInit must be called
// before this.
func NewWorld(name string) (*World, error) {
	w := &World{left: map[string][]*EnemyRuntime{}}
	if err := w.SwitchLevel(name); err != nil {
		return nil, err
	}
	return w, nil
}

// SwitchLevel makes the named level the current one: its platforms, portals,
// background and size. Entering a level brings back its non-boss enemies
// fresh at their spawns (D4); bosses of a level visited before stay as the
// player left them, killed ones dead. The player is not moved.
func (w *World) SwitchLevel(name string) error {
	def, ok := Levels[name]
	if !ok {
		return fmt.Errorf("unknown level %q", name)
	}
	if w.Name != "" {
		w.left[w.Name] = w.Enemies
	}

	level := LoadLevel(def.File, def.Key)
	for i := range level.Portals {
		level.Portals[i].Level = name
	}
	qt := NewDynamicQuadtree(AABB{Width: def.Width, Height: def.Height})
	for i := range level.Platforms {
		qt.Insert(&level.Platforms[i])
	}
	w.border = levelBorder(max(def.Width, level.Width), max(def.Height, level.Height))
	for i := range w.border {
		qt.Insert(&w.border[i])
	}

	enemies, visited := w.left[name]
	if visited {
		enemies = RespawnEnemies(level.Enemies, enemies)
	} else {
		enemies = SpawnEnemies(level.Enemies)
	}
	delete(w.left, name)

	w.Name, w.Def, w.Level, w.Quadtree, w.Enemies = name, def, level, qt, enemies
	return nil
}

// Travel takes p to the portal ref, standing in it, switching level first
// when the portal is in another one. That portal becomes the checkpoint (D10)
// and the level's non-boss enemies revive, also inside the same level (D11).
func (w *World) Travel(p *PlayerRuntime, ref PortalRef) error {
	if ref.Level != w.Name {
		if err := w.SwitchLevel(ref.Level); err != nil {
			return err
		}
	}
	for _, c := range w.Level.Portals {
		if c.Ref() == ref {
			p.placeAt(c.StandPos())
			p.activatePortal(c)
			w.Enemies = RespawnEnemies(w.Level.Enemies, w.Enemies)
			return nil
		}
	}
	return fmt.Errorf("no portal at %v in level %q", ref.Cell, ref.Level)
}

// OpenNextLevel takes p to the first portal of the level after the current
// one in LevelOrder, which becomes used and the checkpoint (levels decision
// D5). It reports false when the current level is the last.
func (w *World) OpenNextLevel(p *PlayerRuntime) (bool, error) {
	i := slices.Index(LevelOrder, w.Name)
	if i < 0 || i+1 >= len(LevelOrder) {
		return false, nil
	}
	next := LevelOrder[i+1]
	if err := w.SwitchLevel(next); err != nil {
		return false, err
	}
	if len(w.Level.Portals) == 0 {
		return false, fmt.Errorf("level %q has no portal to arrive at", next)
	}
	return true, w.Travel(p, w.Level.Portals[0].Ref())
}

// Respawned follows a player respawn: back to the checkpoint's level when it
// is another one (respawn placed p at RespawnPos), and the level's non-boss
// enemies revived.
func (w *World) Respawned(p *PlayerRuntime) {
	if p.Checkpoint.Level != "" && p.Checkpoint.Level != w.Name {
		if err := w.Travel(p, p.Checkpoint); err != nil {
			log.Printf("respawn: %v", err)
		}
	}
	w.Enemies = RespawnEnemies(w.Level.Enemies, w.Enemies)
}
