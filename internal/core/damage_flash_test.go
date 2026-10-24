package core

import "testing"

// tickFlash runs the combat update (which counts the flash down) for n ticks
// with no world, so no hazards or portals interfere.
func tickFlash(p *PlayerRuntime, n int) {
	for i := 0; i < n; i++ {
		p.updateCombat(nil, nil, true)
	}
}

func TestDamageFlashBlinksRedForOneSecond(t *testing.T) {
	p := staminaPlayer()
	if p.flashingRed() {
		t.Fatal("undamaged player is red")
	}

	p.ApplyDamage(5)
	if !p.flashingRed() {
		t.Fatal("not red right after a hit")
	}

	half := int(tps() * damageBlinkSeconds) // one blink half
	tickFlash(p, half+1)
	if p.flashingRed() {
		t.Fatalf("still red after %d ticks, want the normal half of the blink", half+1)
	}
	tickFlash(p, half)
	if !p.flashingRed() {
		t.Fatal("blink did not return to red")
	}

	tickFlash(p, int(tps()*DamageFlashSeconds))
	if p.flashingRed() || p.damageFlash != 0 {
		t.Fatalf("flash left after 1s: %v", p.damageFlash)
	}
}

func TestZeroDamageDoesNotFlash(t *testing.T) {
	p := staminaPlayer()
	p.ApplyDamage(0)
	if p.flashingRed() {
		t.Error("zero damage started the flash")
	}
}

func TestRespawnClearsDamageFlash(t *testing.T) {
	p := staminaPlayer()
	p.ApplyDamage(p.Combat.Health) // lethal
	tickFlash(p, 1)                // updateCombat starts the death animation
	p.respawn()                    // what the end of that animation does
	if p.Combat.Health != p.Combat.MaxHealth || p.flashingRed() {
		t.Errorf("after respawn: health=%d red=%v, want full health and no flash", p.Combat.Health, p.flashingRed())
	}
}
