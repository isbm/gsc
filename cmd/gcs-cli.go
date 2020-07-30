package main

import (
	"fmt"
	"os"

	gsc_clone "github.com/isbm/gsc/clone"

	"github.com/urfave/cli/v2"
)

// Clone package and sync with the git
func clone(ctx *cli.Context) error {
	if ctx.Args().Len() < 2 || ctx.Args().Len() > 3 {
		return fmt.Errorf("Two or three arguments are expected: <project> <name> [git repo]")
	}
	clone := gsc_clone.NewGCSClone().
		SetProject(ctx.Args().Get(0)).
		SetPackage(ctx.Args().Get(1))

	if ctx.Args().Len() == 3 {
		clone.SetGitRepoUrl(ctx.Args().Get(2))
	}

	return clone.Clone()
}

func notImplemented(ctx *cli.Context) error {
	return fmt.Errorf("This feature is not yet implemented, sorry.\nBut you can always send your PR! :)\n")
}

func main() {
	appname := "gsc"
	app := &cli.App{
		Version: "0.1 Alpha",
		Name:    appname,
		Usage:   "OSC to Git binder",
	}

	app.Commands = []*cli.Command{
		{
			Name:   "clone",
			Action: clone,
			Usage:  "<project> <name>\n\tExample: my:cool:project my_package",
		},
		{
			Name:    "checkout",
			Aliases: []string{"co", "bco"},
			Action:  notImplemented,
		},
	}
	if err := app.Run(os.Args); err != nil {
		fmt.Println("Error:", err.Error())
	}
}
