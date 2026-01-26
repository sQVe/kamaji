package main

import "errors"

// configFile is the standard filename for kamaji configuration.
const configFile = "kamaji.yaml"

// Sentinel errors for CLI commands.
// When returned, the error message has already been printed via output package.
var (
	errConfigInvalid = errors.New("config invalid")
	errFileExists    = errors.New("file exists")
	errWriteFailed   = errors.New("write failed")
	errSprintFailed  = errors.New("sprint failed")
)
