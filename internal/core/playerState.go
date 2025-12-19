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
	PlayerStateInAir
	PlayerStateDamaged
	PlayerStateDead
	// ---------------- menu ----------------
	PlayerStateMenuOpen
	// ---------------- skills ----------------
	PlayerStateWeakAttack
	PlayerStateStrongAttack
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

func (ps *PlayerState) IsRunning() bool {
	return ps.CurrentState == PlayerStateRunning
}

// currently we can not able to jump while attacking or defending
func (ps *PlayerState) IsGrounded() bool {
	return ps.CurrentState == PlayerStateIdle || ps.CurrentState == PlayerStateMoving || ps.CurrentState == PlayerStateRunning
}

func (ps *PlayerState) IsInAir() bool {
	return ps.CurrentState == PlayerStateInAir
}

func (ps *PlayerState) IsDamaged() bool {
	return ps.CurrentState == PlayerStateDamaged
}

func (ps *PlayerState) IsDead() bool {
	return ps.CurrentState == PlayerStateDead
}

func (ps *PlayerState) IsWeakAttack() bool {
	return ps.CurrentState == PlayerStateWeakAttack
}

func (ps *PlayerState) IsStrongAttack() bool {
	return ps.CurrentState == PlayerStateStrongAttack
}

func (ps *PlayerState) IsDefense() bool {
	return ps.CurrentState == PlayerStateDefense
}

func (ps *PlayerState) IsUsingPotion() bool {
	return ps.CurrentState == PlayerStateUsePotion
}

func (ps *PlayerState) GetPlayerState() int {
	return int(ps.CurrentState)
}

func (ps *PlayerState) SetPlayerState(state int) {
	ps.CurrentState = PlayerStateType(state)
}
