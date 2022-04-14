//go:build ignore

// This file exists to prevent 'go mod tidy' removing goapi-gen dependancy
// which makes you unable to regenerate the api until added to go.mod again
// It doesn't get built into the final binary.
package api

import (
	_ "github.com/discord-gophers/goapi-gen"
)
