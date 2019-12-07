package version

import (
	"fmt"
	"runtime"

	"github.com/alecthomas/kingpin"
	"github.com/gosuri/uitable"
)

var (
	// GitVersion overridden at build time by:
	//   -ldflags='-X github.com/codedropau/rig/cmd/rig/version.GitVersion=$(git describe --tags --always)'
	GitVersion string
	// GitCommit overridden at build time by:
	//   -ldflags='-X github.com/codedropau/rig/cmd/rig/version.GitCommit=$(git rev-list -1 HEAD)'
	GitCommit string
)

type cmdVersion struct {
	APICompatibility int
	BuildDate        string
	BuildVersion     string
	GOARCH           string
	GOOS             string
}

func (cmd *cmdVersion) run(c *kingpin.ParseContext) error {
	table := uitable.New()
	table.MaxColWidth = 80
	table.AddRow("Version:", GitVersion)
	table.AddRow("Commit:", GitCommit)
	table.AddRow("OS:", runtime.GOOS)
	table.AddRow("Arch:", runtime.GOARCH)
	fmt.Println(table)
	return nil
}

// Command declares the "version" sub command.
func Command(app *kingpin.Application) {
	cmd := new(cmdVersion)
	app.Command("version", fmt.Sprintf("Prints %s version", app.Name)).Action(cmd.run)
}
