//go:build !js

package system

// IsTouchDevice reports false on desktop builds; the on-screen controls still
// turn on there at the first touch (see core.TouchControls.Update).
func IsTouchDevice() bool { return false }
