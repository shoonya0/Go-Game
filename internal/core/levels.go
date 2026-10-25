package core

import (
	"fmt"
	"log"
)

// LevelDef says where a named level comes from: the JSON bundle and key its
// grid is read from, its background image and its size in world pixels.
type LevelDef struct {
	File, Key     string
	Background    string
	Width, Height float64
}

// StartLevel is the level a new game begins in.
const StartLevel = "deepForest1"

// Levels lists every level by name. Travel between portals switches by name.
var Levels = map[string]LevelDef{
	StartLevel: {
		File: LevelFile, Key: DefaultLevelName, Background: Background_1,
		Width: Level_1_Width, Height: Level_1_Height,
	},
}

// World is the level the player is in: its platforms (indexed in Quadtree)
// and its enemies. It remembers the enemies of levels the player left.
type World struct {
	Name     string
	Def      LevelDef
	Level    Level
	Quadtree *DynamicQuadtree
	Enemies  []*EnemyRuntime

	left map[string][]*EnemyRuntime // enemies of each level left, as they were (for its bosses)
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
// when the portal is in another one.
func (w *World) Travel(p *PlayerRuntime, ref PortalRef) error {
	if ref.Level != w.Name {
		if err := w.SwitchLevel(ref.Level); err != nil {
			return err
		}
	}
	for _, c := range w.Level.Portals {
		if c.Ref() == ref {
			p.placeAt(c.StandPos())
			return nil
		}
	}
	return fmt.Errorf("no portal at %v in level %q", ref.Cell, ref.Level)
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
