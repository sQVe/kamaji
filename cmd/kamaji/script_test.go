//go:build integration

package main

import (
	"testing"

	"github.com/rogpeppe/go-internal/testscript"
)

func TestScript(t *testing.T) {
	testscript.Run(t, testscript.Params{
		Dir: "testdata/script",
	})
}

func TestMain(m *testing.M) {
	testscript.Main(m, map[string]func(){
		"kamaji": main,
	})
}
