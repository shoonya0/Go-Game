package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"slices"
)

// saveVersion is written into every save; a save with another version is
// not read (main-menu D2).
const saveVersion = 1

// The storage keys of the two records the game keeps (main-menu D2): the
// progress, wiped by New Game, and the settings, kept.
const (
	SaveKey     = "save"
	SettingsKey = "settings"
)

// SaveData is the progress kept across restarts: the portals the player has
// used and the one they respawn at. Continue starts standing in Checkpoint.
type SaveData struct {
	Version       int         `json:"version"`
	Checkpoint    PortalRef   `json:"checkpoint"`
	ActivePortals []PortalRef `json:"activePortals"`
}

// Snapshot returns p's progress as a save.
func (p *PlayerRuntime) Snapshot() SaveData {
	return SaveData{Version: saveVersion, Checkpoint: p.Checkpoint, ActivePortals: slices.Clone(p.ActivePortals)}
}

// EncodeSave turns s into the bytes the storage keeps.
func EncodeSave(s SaveData) ([]byte, error) { return json.Marshal(s) }

// DecodeSave reads a save written by EncodeSave. A save of another version,
// or without a checkpoint, is an error.
func DecodeSave(b []byte) (SaveData, error) {
	var s SaveData
	if err := json.Unmarshal(b, &s); err != nil {
		return SaveData{}, fmt.Errorf("read save: %w", err)
	}
	if s.Version != saveVersion {
		return SaveData{}, fmt.Errorf("save version %d, want %d", s.Version, saveVersion)
	}
	if s.Checkpoint.Level == "" {
		return SaveData{}, errors.New("save has no checkpoint")
	}
	return s, nil
}

// Continue puts p back where save s left off: standing in its checkpoint
// (switching to that level), with the saved portals active again. Portals that
// no longer exist in their level are dropped; a missing checkpoint is an error
// and changes nothing.
func (w *World) Continue(p *PlayerRuntime, s SaveData) error {
	levels := map[string]Level{}
	exists := func(ref PortalRef) bool {
		if _, ok := Levels[ref.Level]; !ok {
			return false
		}
		level, ok := levels[ref.Level]
		if !ok {
			def := Levels[ref.Level]
			level = LoadLevel(def.File, def.Key)
			levels[ref.Level] = level
		}
		return slices.ContainsFunc(level.Portals, func(c Portal) bool { return c.Pos == ref.Cell })
	}
	if !exists(s.Checkpoint) {
		return fmt.Errorf("saved checkpoint %v in level %q does not exist", s.Checkpoint.Cell, s.Checkpoint.Level)
	}
	active := make([]PortalRef, 0, len(s.ActivePortals))
	for _, ref := range s.ActivePortals {
		if exists(ref) && !slices.Contains(active, ref) {
			active = append(active, ref)
		}
	}
	if err := w.Travel(p, s.Checkpoint); err != nil {
		return err
	}
	p.ActivePortals = active
	if !slices.Contains(p.ActivePortals, s.Checkpoint) {
		p.ActivePortals = append(p.ActivePortals, s.Checkpoint)
	}
	return nil
}

// Autosave tells when the player's progress differs from the last one
// written, so the game saves on a change only (main-menu D3).
type Autosave struct {
	last    SaveData
	written bool
}

// Due returns p's progress and true when it differs from the last one marked
// written. The caller stores it, then calls Written.
func (a *Autosave) Due(p *PlayerRuntime) (SaveData, bool) {
	s := p.Snapshot()
	if a.written && s.Checkpoint == a.last.Checkpoint && slices.Equal(s.ActivePortals, a.last.ActivePortals) {
		return SaveData{}, false
	}
	return s, true
}

// Written marks s as the progress in storage.
func (a *Autosave) Written(s SaveData) { a.last, a.written = s, true }
