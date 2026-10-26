//go:build js

package system

import (
	"fmt"
	"syscall/js"
)

// FetchFile downloads url (relative to the game's page) with the browser's
// fetch. It blocks until the download ends, so call it from a goroutine, not
// the game loop.
func FetchFile(url string) ([]byte, error) {
	type result struct {
		data []byte
		err  error
	}
	done := make(chan result, 1)
	finish := func(r result) {
		select { // only the first outcome counts; a second must not block JS
		case done <- r:
		default:
		}
	}

	onError := js.FuncOf(func(_ js.Value, args []js.Value) any {
		finish(result{err: fmt.Errorf("fetch %s: %s", url, args[0].Call("toString").String())})
		return nil
	})
	onBody := js.FuncOf(func(_ js.Value, args []js.Value) any {
		bytes := js.Global().Get("Uint8Array").New(args[0])
		data := make([]byte, bytes.Get("length").Int())
		js.CopyBytesToGo(data, bytes)
		finish(result{data: data})
		return nil
	})
	onResponse := js.FuncOf(func(_ js.Value, args []js.Value) any {
		resp := args[0]
		if !resp.Get("ok").Bool() {
			finish(result{err: fmt.Errorf("fetch %s: HTTP %d", url, resp.Get("status").Int())})
			return nil
		}
		resp.Call("arrayBuffer").Call("then", onBody).Call("catch", onError)
		return nil
	})
	defer onError.Release()
	defer onBody.Release()
	defer onResponse.Release()

	js.Global().Call("fetch", url).Call("then", onResponse).Call("catch", onError)
	r := <-done
	return r.data, r.err
}
