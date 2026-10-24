package core

import (
	"bytes"
	"fmt"
	"image"
	"maps"
	"reflect"
	"slices"
	"strings"
	"testing"
)

// TestEnemyMissingClipFallsBack covers a species with only the required clips:
// optional clips play a related one, and guard is simply unavailable.
func TestEnemyMissingClipFallsBack(t *testing.T) {
	bp := succubusBlueprint // copy; its Clips map is replaced below, not shared
	bp.Clips = map[ClipName]EnemyClip{}
	for _, name := range requiredClips {
		bp.Clips[name] = succubusBlueprint.Clips[name]
	}

	for name, want := range map[ClipName]ClipName{
		ClipRun: ClipWalk, ClipAlert: ClipWalk, ClipAttack2: ClipAttack1, ClipIdle: ClipIdle,
	} {
		if got, ok := bp.resolveClip(name); !ok || got != want {
			t.Errorf("resolveClip(%v) = %v, %v; want %v, true", name, got, ok, want)
		}
	}
	if _, ok := bp.resolveClip(ClipGuard); ok {
		t.Error("resolveClip(guard) ok without guard art, want unavailable")
	}

	e := &EnemyRuntime{bp: &bp, Animations: bp.newAnimations(), State: EnemyAlert}
	if e.clipFor(EnemyAlert) != ClipWalk || e.currentAnim() != e.Animations[ClipWalk] {
		t.Errorf("chasing without run or alert art plays %v, want walk", e.clipFor(EnemyAlert))
	}
}

// TestBlueprintHasOwnArtOnly checks Has counts only a species' own art: a
// clip it can play through a fallback is not a capability.
func TestBlueprintHasOwnArtOnly(t *testing.T) {
	bp := succubusBlueprint // copy; its Clips map is replaced below, not shared
	bp.Clips = map[ClipName]EnemyClip{}
	for _, name := range requiredClips {
		bp.Clips[name] = succubusBlueprint.Clips[name]
	}

	for _, name := range requiredClips {
		if !bp.Has(name) {
			t.Errorf("Has(%v) = false for a required clip it has", name)
		}
	}
	for _, name := range []ClipName{ClipAlert, ClipRun, ClipAttack2, ClipGuard} {
		if bp.Has(name) {
			t.Errorf("Has(%v) = true without its art", name)
		}
	}

	want := []string{"idle", "walk", "(walk)", "(walk)", "attack_1", "(attack_1)", "-", "hurt", "death"}
	if got := bp.capabilities(); !reflect.DeepEqual(got, want) {
		t.Errorf("capabilities() = %q\nwant %q", got, want)
	}
}

// TestEnemyCapabilities prints what every registered enemy can do, one row per
// species: run it with -v to see the table.
//
//	go test ./internal/core -run TestEnemyCapabilities -v
func TestEnemyCapabilities(t *testing.T) {
	var b strings.Builder
	fmt.Fprintf(&b, "%-16s", "enemy")
	for _, name := range clipNames {
		fmt.Fprintf(&b, " %-12s", name)
	}
	kinds := slices.Sorted(maps.Keys(enemyBlueprints))
	for _, kind := range kinds {
		bp := enemyBlueprints[kind]
		cells := bp.capabilities()
		if len(cells) != len(clipNames) {
			t.Fatalf("%s: %d cells, want one per clip (%d)", bp.Name, len(cells), len(clipNames))
		}
		fmt.Fprintf(&b, "\n%-16s", bp.Name)
		for _, c := range cells {
			fmt.Fprintf(&b, " %-12s", c)
		}
	}
	t.Log("\n" + b.String())
}

// The succubus has no run clip, so it chases with its alert loop, as before
// clips had standard names.
func TestSuccubusChasesWithAlertClip(t *testing.T) {
	e := newAnimatedTestEnemy(Position{})
	if got := e.clipFor(EnemyAlert); got != ClipAlert {
		t.Errorf("succubus chase clip = %v, want alert", got)
	}
}

func TestEveryEnemyStateHasAClip(t *testing.T) {
	for st := EnemyIdle; st <= EnemyDead; st++ {
		if _, ok := stateClips[st]; !ok {
			t.Errorf("state %d has no clip in stateClips", st)
		}
	}
}

// TestEnemyClipsFitTheirSheets checks every clip's frames lie inside its sheet.
// A frame past the edge draws nothing, so the enemy blinks out for that frame:
// the succubus patrol clip once ran one column past its 8-column row.
func TestEnemyClipsFitTheirSheets(t *testing.T) {
	for _, bp := range enemyBlueprints {
		for name, c := range bp.Clips {
			cfg, _, err := image.DecodeConfig(bytes.NewReader(readAsset(c.Sheet)))
			if err != nil {
				t.Errorf("%s clip %v: decode %q: %v", bp.Name, name, c.Sheet, err)
				continue
			}
			if c.Frames <= 0 || c.FrameW <= 0 || c.FrameH <= 0 || c.FPS <= 0 {
				t.Errorf("%s clip %v: frames=%d size=%dx%d fps=%v, want all positive", bp.Name, name, c.Frames, c.FrameW, c.FrameH, c.FPS)
				continue
			}
			if right := (c.Start + c.Frames) * c.FrameW; right > cfg.Width {
				t.Errorf("%s clip %v: frames end at x=%d, past the %dpx-wide sheet %q", bp.Name, name, right, cfg.Width, c.Sheet)
			}
			if c.HitFrame < 0 || c.HitFrame >= c.Frames {
				t.Errorf("%s clip %v: HitFrame %d outside its %d frames", bp.Name, name, c.HitFrame, c.Frames)
			}
			if len(c.Hold) > c.Frames {
				t.Errorf("%s clip %v: %d Hold entries for %d frames", bp.Name, name, len(c.Hold), c.Frames)
			}
			if bottom := (c.Row + 1) * c.FrameH; bottom > cfg.Height {
				t.Errorf("%s clip %v: row ends at y=%d, past the %dpx-high sheet %q", bp.Name, name, bottom, cfg.Height, c.Sheet)
			}
		}
	}
}

// TestSkeletonVariantsShareBlueprint checks the skeleton colours are the same
// enemy: only name, marker and sheet paths differ.
func TestSkeletonVariantsShareBlueprint(t *testing.T) {
	white, yellow := blueprintFor(KindSkeletonWhite), blueprintFor(KindSkeletonYellow)
	if white.Marker == yellow.Marker || white.Name == yellow.Name {
		t.Errorf("skeleton colours share name or marker: %q %q / %q %q", white.Name, white.Marker, yellow.Name, yellow.Marker)
	}
	if white.Stats != yellow.Stats || white.BodyWidth != yellow.BodyWidth ||
		white.BodyHeight != yellow.BodyHeight || white.Scale != yellow.Scale {
		t.Error("skeleton colours differ in stats, body or scale")
	}
	if len(white.Clips) != len(yellow.Clips) {
		t.Fatalf("skeleton colours have %d and %d clips", len(white.Clips), len(yellow.Clips))
	}
	for name, w := range white.Clips {
		y := yellow.Clips[name]
		if w.Sheet == y.Sheet {
			t.Errorf("clip %v: both colours draw %q", name, w.Sheet)
		}
		w.Sheet, y.Sheet = "", ""
		if !reflect.DeepEqual(w, y) {
			t.Errorf("clip %v differs between colours beyond its sheet: %+v / %+v", name, w, y)
		}
	}
}
