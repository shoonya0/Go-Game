package core

import (
	"math"
	"testing"
)

// TestClipHoldTiming checks a held frame stays up Hold times as long as the
// others, and that frames past the end of Hold keep the base duration.
func TestClipHoldTiming(t *testing.T) {
	bp := succubusBlueprint // copy; its Clips map is replaced below, not shared
	bp.Clips = map[ClipName]EnemyClip{}
	for name, c := range succubusBlueprint.Clips {
		bp.Clips[name] = c
	}
	bp.Clips[ClipWalk] = succubusClip(1, 1, 3, 6, true).held(1, 3) // frame 2 has no entry
	start := bp.Clips[ClipWalk].Start
	e := &EnemyRuntime{bp: &bp, Animations: bp.newAnimations(), State: EnemyPatrol, PrevState: EnemyPatrol, CurrAnimFrame: start}

	base := tps() / 6 // ticks per unheld frame
	want := []float64{base, 3 * base, base}
	ticks := make([]int, len(want))
	for range int(5*base) + 2 { // one full play (1+3+1 frames) and the wrap
		ticks[e.CurrAnimFrame-start]++
		e.UpdateEnemyAnimation()
	}
	for i, w := range want {
		if math.Abs(float64(ticks[i])-w) > 1 {
			t.Errorf("frame %d shown %d ticks, want %v (±1)", i, ticks[i], w)
		}
	}
	if e.CurrAnimFrame != start {
		t.Errorf("after one play frame = %d, want looped back to %d", e.CurrAnimFrame, start)
	}
}

func TestClipDurationCountsHolds(t *testing.T) {
	c := EnemyClip{Frames: 4, FPS: 10}.held(2, 0, 3) // 0 means 1; frame 3 has no entry
	if got, want := c.duration(), (2+1+3+1)/10.0; math.Abs(got-want) > 1e-9 {
		t.Errorf("duration = %v, want %v", got, want)
	}
}

func TestMotionPose(t *testing.T) {
	if p := (motionFX{}).pose(0.3, false); p != (spritePose{}) {
		t.Errorf("zero motion pose = %+v, want none", p)
	}

	bob := motionFX{BobPx: 4, BobHz: 2}
	if p := bob.pose(0.25, false); math.Abs(p.DY+4) > 1e-9 { // top of the first bounce
		t.Errorf("bob at its top: DY = %v, want -4", p.DY)
	}
	for _, tm := range []float64{0, 0.1, 0.37, 0.5, 1.9} {
		if p := bob.pose(tm, false); p.DY > 0 || p.DY < -4 {
			t.Errorf("bob at t=%v: DY = %v, want within [-4, 0] (never below the feet)", tm, p.DY)
		}
	}

	lean := motionFX{LeanDeg: 10}
	right, left := lean.pose(0, false).Rot, lean.pose(0, true).Rot
	if right <= 0 || left != -right {
		t.Errorf("lean rot right=%v left=%v, want clockwise facing right, mirrored facing left", right, left)
	}

	shake := motionFX{HitShakePx: 5}
	if p := shake.pose(1.0/(4*hitShakeHz), false); p.DX <= 0 || p.DX > 5 {
		t.Errorf("shake near its first peak: DX = %v, want in (0, 5]", p.DX)
	}
	if p := shake.pose(hitShakeDuration, false); p.DX != 0 {
		t.Errorf("shake after it died out: DX = %v, want 0", p.DX)
	}
}

// TestEnemyClipTimeRestarts checks the motion clock starts over with each clip,
// so a hit shake plays from its start on every hit.
func TestEnemyClipTimeRestarts(t *testing.T) {
	e := newAnimatedTestEnemy(Position{})
	tickEnemyAnim(e, 30)
	if e.clipTime <= 0 {
		t.Fatalf("clipTime = %v after 30 ticks, want it advancing", e.clipTime)
	}
	e.TakeDamage(1)
	if e.clipTime != 0 {
		t.Errorf("clipTime = %v right after a hit, want 0", e.clipTime)
	}
	e.PrevState = e.State
	tickEnemyAnim(e, clipTicks(e, EnemyHurt))
	if e.State != EnemyPatrol || e.clipTime > 2*deltaTime() {
		t.Errorf("after hurt: state=%d clipTime=%v, want patrol with a fresh clock", e.State, e.clipTime)
	}
}
