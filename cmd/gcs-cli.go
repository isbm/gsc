package main

import (
	"fmt"
	"os"

	gsc_add "github.com/isbm/gsc/add"
	gsc_clone "github.com/isbm/gsc/clone"
	gsc_push "github.com/isbm/gsc/push"
	gsc_submit "github.com/isbm/gsc/submit"
	"github.com/urfave/cli/v2"
)

// Submit request
func submit(ctx *cli.Context) error {
	sr := gsc_submit.NewGSCSubmitRequest()
	sr.Submit()
	return nil
}

// Push package to the git and OBS
func push(ctx *cli.Context) error {
	push := gsc_push.NewGCSPush()
	return push.Push()
}

func add(ctx *cli.Context) error {
	add := gsc_add.NewGSCAdd()
	if ctx.Args().Len() > 0 {
		add.SetPathspec(ctx.Args().Get(0))
	}
	return add.Add()
}

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
			Usage:  "Clone package from the OBS and link to Git repo.\n\tUsage: <project> <name>\n\tExample: my:cool:project my_package [repo.git]",
		},
		{
			Name:    "add",
			Aliases: []string{"a"},
			Action:  add,
		},
		{
			Name:    "submitreq",
			Aliases: []string{"sr"},
			Action:  submit,
			Usage:   "Create request to submit source back to Project",
		},
		{
			Name:    "push",
			Aliases: []string{"p"},
			Action:  notImplemented,
			Usage:   "Push package to the OBS branch and Git repo",
		},
	}
	if err := app.Run(os.Args); err != nil {
		fmt.Println("Error:", err.Error())
	}
}
