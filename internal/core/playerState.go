package core

type PlayerStateType int

// ---------------- player state ----------------
const (
	// ---------------- states ----------------
	PlayerStateIdle PlayerStateType = iota
	PlayerStateMoving
	PlayerStateRunning
	PlayerStateJumping
	PlayerStateFalling
	PlayerStateLanding
	PlayerStateSmugFace
	PlayerStateWeakAttack
	PlayerStateStrongAttack
	PlayerStateSpecialAttack1
	PlayerStateSpecialAttack2
	PlayerStateSpecialAttack3
	PlayerStateSpecialAttack4
	PlayerStateWeakAttackInAir
	PlayerStateStrongAttackInAir

	// not used yet
	PlayerStateDamaged
	PlayerStateDead
	// ---------------- menu ----------------
	PlayerStateMenuOpen
	// ---------------- skills ----------------
	PlayerStateDefense
	PlayerStateUsePotion
)

type PlayerState struct {
	CurrentState PlayerStateType
}

func (ps *PlayerState) IsIdle() bool {
	return ps.CurrentState == PlayerStateIdle
}

func (ps *PlayerState) IsMoving() bool {
	return ps.CurrentState == PlayerStateMoving
}

func (ps *PlayerState) IsJumping() bool {
	return ps.CurrentState == PlayerStateJumping
}

func (ps *PlayerState) IsFalling() bool {
	return ps.CurrentState == PlayerStateFalling
}

func (ps *PlayerState) IsLanding() bool {
	return ps.CurrentState == PlayerStateLanding
}

func (ps *PlayerState) IsRunning() bool {
	return ps.CurrentState == PlayerStateRunning
}

func (ps *PlayerState) IsSmugFace() bool {
	return ps.CurrentState == PlayerStateSmugFace
}

func (ps *PlayerState) IsWeakAttack() bool {
	return ps.CurrentState == PlayerStateWeakAttack
}

func (ps *PlayerState) IsStrongAttack() bool {
	return ps.CurrentState == PlayerStateStrongAttack
}

func (ps *PlayerState) IsSpecialAttack1() bool {
	return ps.CurrentState == PlayerStateSpecialAttack1
}

func (ps *PlayerState) IsSpecialAttack2() bool {
	return ps.CurrentState == PlayerStateSpecialAttack2
}

func (ps *PlayerState) IsSpecialAttack3() bool {
	return ps.CurrentState == PlayerStateSpecialAttack3
}

func (ps *PlayerState) IsSpecialAttack4() bool {
	return ps.CurrentState == PlayerStateSpecialAttack4
}

func (ps *PlayerState) IsWeakAttackInAir() bool {
	return ps.CurrentState == PlayerStateWeakAttackInAir
}

func (ps *PlayerState) IsStrongAttackInAir() bool {
	return ps.CurrentState == PlayerStateStrongAttackInAir
}

// currently we can not able to jump while attacking or defending
func (ps *PlayerState) IsGrounded() bool {
	return ps.CurrentState == PlayerStateIdle || ps.CurrentState == PlayerStateMoving || ps.CurrentState == PlayerStateRunning
}

func (ps *PlayerState) IsDamaged() bool {
	return ps.CurrentState == PlayerStateDamaged
}

func (ps *PlayerState) IsDead() bool {
	return ps.CurrentState == PlayerStateDead
}

func (ps *PlayerState) IsDefense() bool {
	return ps.CurrentState == PlayerStateDefense
}

func (ps *PlayerState) IsUsingPotion() bool {
	return ps.CurrentState == PlayerStateUsePotion
}

// IsGroundedOneShot reports whether the current state is a non-looping grounded
// action that should return to idle once its animation finishes.
func (ps *PlayerState) IsGroundedOneShot() bool {
	return ps.IsLanding() || ps.IsSmugFace() || ps.IsWeakAttack() || ps.IsStrongAttack() ||
		ps.IsSpecialAttack1() || ps.IsSpecialAttack2() || ps.IsSpecialAttack3() || ps.IsSpecialAttack4()
}

// IsAirOneShot reports whether the current state is a non-looping airborne
// action that should transition to falling once its animation finishes.
func (ps *PlayerState) IsAirOneShot() bool {
	return ps.IsWeakAttackInAir() || ps.IsStrongAttackInAir() || ps.IsJumping()
}

func (ps *PlayerState) GetPlayerState() int {
	return int(ps.CurrentState)
}

func (ps *PlayerState) SetPlayerState(state int) {
	ps.CurrentState = PlayerStateType(state)
}
