//go:build !js

package system

import "errors"

// FetchFile is for the web build only: desktop builds embed every asset.
func FetchFile(url string) ([]byte, error) {
	return nil, errors.New("fetch " + url + ": only the web build downloads assets")
}
