package main

import (
	"os"

	"github.com/alecthomas/kingpin"

	"github.com/nickschuch/rig/cmd/rig/run"
	"github.com/nickschuch/rig/cmd/rig/snapshot"
)

func main() {
	app := kingpin.New("rig", "Docker Compose cloning tool")

	snapshot.Command(app)
	run.Command(app)

	kingpin.MustParse(app.Parse(os.Args[1:]))
}