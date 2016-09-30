package main

import (
	"os"

	"github.com/elastic/beats/libbeat/beat"

	"github.com/mod/loxybeat/beater"
)

func main() {
	err := beat.Run("loxybeat", "", beater.New)
	if err != nil {
		os.Exit(1)
	}
}
