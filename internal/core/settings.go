package core

import (
	"encoding/json"
	"fmt"
)

// MaxVolume is the loudest sound volume; volume steps by 1 from 0 (main-menu D5).
const MaxVolume = 10

// Settings are the player's choices, kept apart from the save so New Game
// keeps them (D2). No audio plays yet; Volume is stored for when it does (D5).
type Settings struct {
	Volume int `json:"volume"`
}

// DefaultSettings are the settings before the player changes any.
func DefaultSettings() Settings { return Settings{Volume: 7} }

// EncodeSettings turns s into the bytes the storage keeps.
func EncodeSettings(s Settings) ([]byte, error) { return json.Marshal(s) }

// DecodeSettings reads settings written by EncodeSettings; a volume out of
// range is an error.
func DecodeSettings(b []byte) (Settings, error) {
	s := DefaultSettings()
	if err := json.Unmarshal(b, &s); err != nil {
		return DefaultSettings(), fmt.Errorf("read settings: %w", err)
	}
	if s.Volume < 0 || s.Volume > MaxVolume {
		return DefaultSettings(), fmt.Errorf("volume %d out of 0..%d", s.Volume, MaxVolume)
	}
	return s, nil
}

// volumeLabel is the settings screen's volume item: arrows say it changes
// with left/right.
func volumeLabel(v int) string { return fmt.Sprintf("< Volume %2d >", v) }
