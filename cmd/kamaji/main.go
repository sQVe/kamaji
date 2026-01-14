package main

import (
	"fmt"

	"github.com/sqve/kamaji/internal/version"
)

func main() {
	fmt.Println("kamaji", version.Full())
}
