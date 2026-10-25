//go:build !js

package system

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
)

// saveDir returns the folder the records are kept in: the game's folder in
// the user's config folder (%AppData% on Windows). Tests replace it.
var saveDir = func() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "PirateAdventure"), nil
}

func savePath(key string) (string, error) {
	dir, err := saveDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, key+".json"), nil
}

// LoadRecord returns the record stored under key; ok is false when there is
// none or it can't be read.
func LoadRecord(key string) (data []byte, ok bool) {
	path, err := savePath(key)
	if err != nil {
		return nil, false
	}
	data, err = os.ReadFile(path)
	return data, err == nil
}

// StoreRecord keeps data under key, replacing the old record. It writes a
// temporary file first, so a crash mid-write leaves the old record whole.
func StoreRecord(key string, data []byte) error {
	path, err := savePath(key)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}

// DeleteRecord removes the record stored under key; none is not an error.
func DeleteRecord(key string) error {
	path, err := savePath(key)
	if err != nil {
		return err
	}
	if err := os.Remove(path); err != nil && !errors.Is(err, fs.ErrNotExist) {
		return err
	}
	return nil
}
