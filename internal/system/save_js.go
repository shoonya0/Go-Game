//go:build js

package system

import (
	"errors"
	"fmt"
	"syscall/js"
)

// recordPrefix keeps the game's keys apart from other pages on the same origin.
const recordPrefix = "pirateAdventure."

// localStorage returns the browser's storage, or an error when the page can't
// use it (some private windows throw on access).
func localStorage() (store js.Value, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("localStorage: %v", r)
		}
	}()
	store = js.Global().Get("localStorage")
	if !store.Truthy() {
		return js.Value{}, errors.New("localStorage unavailable")
	}
	return store, nil
}

// call runs store.method(args...), turning a thrown JS error (quota full,
// storage blocked) into an error.
func call(store js.Value, method string, args ...any) (v js.Value, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("localStorage.%s: %v", method, r)
		}
	}()
	return store.Call(method, args...), nil
}

// LoadRecord returns the record stored under key; ok is false when there is
// none or storage can't be used.
func LoadRecord(key string) (data []byte, ok bool) {
	store, err := localStorage()
	if err != nil {
		return nil, false
	}
	v, err := call(store, "getItem", recordPrefix+key)
	if err != nil || v.Type() != js.TypeString {
		return nil, false
	}
	return []byte(v.String()), true
}

// StoreRecord keeps data under key, replacing the old record.
func StoreRecord(key string, data []byte) error {
	store, err := localStorage()
	if err != nil {
		return err
	}
	_, err = call(store, "setItem", recordPrefix+key, string(data))
	return err
}

// DeleteRecord removes the record stored under key; none is not an error.
func DeleteRecord(key string) error {
	store, err := localStorage()
	if err != nil {
		return err
	}
	_, err = call(store, "removeItem", recordPrefix+key)
	return err
}
