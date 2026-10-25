package core

// Direction is the discrete axis input for a frame (-1, 0 or 1 per axis).
type Direction struct {
	LeftRight int8
	UpDown    int8
}

// InputState is the fully-decoded player input for a single frame. It is
// produced by the input adapter (see internal/system) and consumed by the
// player update, keeping core free of any device/keyboard dependency.
type InputState struct {
	Direction       Direction
	JumpJustPressed bool     // jump button pressed this frame
	JumpHeld        bool     // jump button currently held
	DashJustPressed bool     // dash button pressed this frame
	RunJustPressed  bool     // run modifier pressed this frame
	Menu            bool     // menu toggled
	SmugFace        bool     // taunt button pressed
	ToggleHitboxes  bool     // debug hitbox overlay toggled this frame
	Interact        bool     // use button (E) pressed this frame: saves the portal the player stands in
	Tapped          bool     // a touch off the on-screen buttons started this frame
	TapAt           Position // that touch's screen point
	Skills          Skills
}
