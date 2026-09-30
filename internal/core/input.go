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
	JumpJustPressed bool // jump button pressed this frame
	JumpHeld        bool // jump button currently held
	DashJustPressed bool // dash button pressed this frame
	RunJustPressed  bool // run modifier pressed this frame
	Menu            bool // menu toggled
	SmugFace        bool // taunt button pressed
	Skills          Skills
}
