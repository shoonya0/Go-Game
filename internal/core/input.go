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
	InteractTicks   int      // ticks the use button (E) has been held, 0 while up: a long hold opens the travel menu
	Tapped          bool     // a touch off the on-screen buttons started this frame
	TouchHoldTicks  int      // ticks that touch has stayed down since it started, 0 when there is none
	TapAt           Position // that touch's screen point (where it is now, while held)
	UpJustPressed   bool     // up pressed this frame (menus)
	DownJustPressed bool     // down pressed this frame (menus)
	Confirm         bool     // Enter pressed this frame (menus)
	Back            bool     // Escape pressed this frame: closes a menu
	Skills          Skills
}
