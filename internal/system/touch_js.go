//go:build js

package system

import "syscall/js"

// IsTouchDevice reports whether the browser has a touch screen (phones and
// tablets), so the on-screen controls show before the first touch.
func IsTouchDevice() bool {
	nav := js.Global().Get("navigator")
	if nav.Truthy() {
		if n := nav.Get("maxTouchPoints"); n.Type() == js.TypeNumber && n.Int() > 0 {
			return true
		}
	}
	return js.Global().Get("ontouchstart").Type() != js.TypeUndefined
}
