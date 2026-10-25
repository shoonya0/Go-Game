package core

import "fmt"

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

	left map[string][]*EnemyRuntime // enemies of each level left, as they were
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
// background and size. A level visited before gets back its enemies as the
// player left them; a new one spawns them fresh. The player is not moved.
func (w *World) SwitchLevel(name string) error {
	def, ok := Levels[name]
	if !ok {
		return fmt.Errorf("unknown level %q", name)
	}
	if w.Name != "" {
		w.left[w.Name] = w.Enemies
	}

	level := LoadLevel(def.File, def.Key)
	qt := NewDynamicQuadtree(AABB{Width: def.Width, Height: def.Height})
	for i := range level.Platforms {
		qt.Insert(&level.Platforms[i])
	}

	enemies, visited := w.left[name]
	if !visited {
		enemies = SpawnEnemies(level.Enemies)
	}
	delete(w.left, name)

	w.Name, w.Def, w.Level, w.Quadtree, w.Enemies = name, def, level, qt, enemies
	return nil
}
