package main

import (
	"os"

	"github.com/alecthomas/kingpin"

	"github.com/codedropau/rig/cmd/rig/list"
	"github.com/codedropau/rig/cmd/rig/run"
	"github.com/codedropau/rig/cmd/rig/snapshot"
	"github.com/codedropau/rig/cmd/rig/version"
)

func main() {
	app := kingpin.New("rig", "Docker Compose cloning tool")

	snapshot.Command(app)
	run.Command(app)
	list.Command(app)
	version.Command(app)

	kingpin.MustParse(app.Parse(os.Args[1:]))
}
