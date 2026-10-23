// Command battler2strip bakes an RPG Maker MV side-view battler sheet into the
// game's strip format: one right-facing PNG strip per clip, shrunk to the size
// drawn in game, plus a blueprint snippet printed to stdout. It runs offline
// and is not part of the game build (enemy-types decision D7).
//
//	go run ./tools/battler2strip -out assets/enemy/Goblin tools/battler2strip/mappings/goblin.json
//
// Run it from the repo root: the mapping's sheet path and -out are relative to
// the working directory.
package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"image"
	"image/png"
	"os"
	"path/filepath"
	"sort"
)

func main() {
	out := flag.String("out", "", "folder to write the strips to (created if missing)")
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "usage: battler2strip -out DIR MAPPING.json")
		flag.PrintDefaults()
	}
	flag.Parse()
	if *out == "" || flag.NArg() != 1 {
		flag.Usage()
		os.Exit(2)
	}
	if err := run(flag.Arg(0), *out); err != nil {
		fmt.Fprintln(os.Stderr, "battler2strip:", err)
		os.Exit(1)
	}
}

func run(mappingPath, outDir string) error {
	m, err := readMapping(mappingPath)
	if err != nil {
		return err
	}
	sheet, err := readPNG(m.Sheet)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return err
	}

	names := make([]string, 0, len(m.Clips))
	for name := range m.Clips {
		names = append(names, name)
	}
	sort.Strings(names)
	var total int
	for _, name := range names {
		strip, err := bakeClip(sheet, m, m.Clips[name])
		if err != nil {
			return fmt.Errorf("clip %q: %w", name, err)
		}
		var buf bytes.Buffer
		if err := png.Encode(&buf, strip); err != nil {
			return err
		}
		file := filepath.Join(outDir, stripName(m, name))
		if err := os.WriteFile(file, buf.Bytes(), 0o644); err != nil {
			return err
		}
		total += buf.Len()
		fmt.Fprintf(os.Stderr, "wrote %s (%d frames, %d KB)\n", file, len(m.Clips[name].Frames), buf.Len()/1024)
	}
	fmt.Fprintf(os.Stderr, "%d strips, %d KB in all\n", len(names), total/1024)

	w, h := frameSize(sheet.Bounds(), m.DrawHeight)
	fmt.Print(snippet(m, repoRelative(outDir), w, h))
	return nil
}

// repoRelative is dir relative to the working directory (the repo root), so
// the snippet's sheet paths look like the other blueprints' whether -out was
// given relative or absolute.
func repoRelative(dir string) string {
	abs, err := filepath.Abs(dir)
	if err != nil {
		return dir
	}
	wd, err := os.Getwd()
	if err != nil {
		return dir
	}
	if rel, err := filepath.Rel(wd, abs); err == nil {
		return rel
	}
	return dir
}

// readMapping decodes and checks a mapping file; unknown fields are errors, so
// a typo doesn't silently drop a setting.
func readMapping(file string) (Mapping, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return Mapping{}, err
	}
	var m Mapping
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&m); err != nil {
		return Mapping{}, fmt.Errorf("%s: %w", file, err)
	}
	if err := m.validate(); err != nil {
		return Mapping{}, fmt.Errorf("%s: %w", file, err)
	}
	return m, nil
}

func readPNG(file string) (image.Image, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", file, err)
	}
	return img, nil
}
