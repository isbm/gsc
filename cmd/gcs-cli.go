package main

import (
	"fmt"
	"os"

	wzlib_logger "github.com/infra-whizz/wzlib/logger"
	gsc_add "github.com/isbm/gsc/add"
	gsc_clone "github.com/isbm/gsc/clone"
	gsc_close "github.com/isbm/gsc/close"
	gsc_import "github.com/isbm/gsc/import"
	gsc_merge "github.com/isbm/gsc/merge_branch"
	gsc_info "github.com/isbm/gsc/pkginfo"
	gsc_release "github.com/isbm/gsc/release"
	gsc_submit "github.com/isbm/gsc/submit"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func setLogger(ctx *cli.Context) {
	if ctx.Bool("debug") {
		wzlib_logger.GetCurrentLogger().SetLevel(logrus.DebugLevel)
	} else {
		wzlib_logger.GetCurrentLogger().SetLevel(logrus.InfoLevel)
	}
}

// Release package to a branch
func releasePackage(ctx *cli.Context) error {
	setLogger(ctx)
	rel := gsc_release.NewGSCPackageRelease()
	rel.SetReleaseBranch(ctx.String("branch"))

	return rel.Release()
}

// Import package (SR was accepted)
func importPackage(ctx *cli.Context) error {
	setLogger(ctx)
	imp := gsc_import.NewGSCPackageImport()
	repo := ctx.String("git-repo")
	if repo != "" {
		imp.SetGitRepoUrl(repo)
	}
	return imp.Import()
}

// Merge package (SR was accepted)
func merge(ctx *cli.Context) error {
	setLogger(ctx)
	return gsc_merge.NewGSCMergeBranch().Merge(ctx.Bool("develop"))
}

// Submit request
func submit(ctx *cli.Context) error {
	setLogger(ctx)
	return gsc_submit.NewGSCSubmitRequest().Submit()
}

// Get package information
func packageInfo(ctx *cli.Context) error {
	setLogger(ctx)
	info := gsc_info.NewGSCPackageInfo()
	if err := info.ObtainInfo(); err != nil {
		return err
	}

	fmt.Printf("       Package: %s\n", info.Name)
	fmt.Printf("       Version: %s\n", info.Version)
	fmt.Printf("      Revision: %s\n", info.Project.Revision)
	fmt.Printf("      Revision: %s\n\n", info.Project.ApiUrl)
	fmt.Printf("Git Repository: %s\n", info.GitRepo.Url)
	fmt.Printf("    Git Branch: %s\n", info.GitRepo.Branch)
	fmt.Printf("    Source URL: %s\n", info.Project.SourceUrl)

	return nil
}

// Close current branch
func closeBranch(ctx *cli.Context) error {
	setLogger(ctx)
	cls := gsc_close.NewGSCCloseBranch()
	if err := cls.Close(); err != nil {
		return err
	}
	return cls.Cleanup()
}

func add(ctx *cli.Context) error {
	setLogger(ctx)
	add := gsc_add.NewGSCAdd()
	if ctx.Args().Len() > 0 {
		add.SetPathspec(ctx.Args().Get(0))
	}
	return add.Add()
}

// Clone package and sync with the git
func clone(ctx *cli.Context) error {
	setLogger(ctx)
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
	app.Flags = []cli.Flag{
		&cli.BoolFlag{
			Name:    "debug",
			Aliases: []string{"d"},
			Usage:   "Turn on debug level",
		},
	}

	app.Commands = []*cli.Command{
		{
			Name:    "clone",
			Aliases: []string{"co"},
			Action:  clone,
			Usage:   "Clone package from the OBS and link to Git repo.\n\tUsage: <project> <name>\n\tExample: my:cool:project my_package [repo.git]",
		},
		{
			Name:    "import",
			Aliases: []string{"i", "im"},
			Action:  importPackage,
			Usage:   "Import package from the sources",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "git-repo",
					Aliases: []string{"gr", "r"},
					Usage:   "Git repository to import package from",
				},
			},
		},
		{
			Name:    "add",
			Aliases: []string{"a"},
			Action:  add,
			Usage:   "Add modified changes to the current session",
		},
		{
			Name:    "release",
			Aliases: []string{"rl", "rel"},
			Action:  releasePackage,
			Usage:   "Release to a branch",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "branch",
					Aliases:  []string{"b"},
					Usage:    "Name of the release branch",
					Required: true,
				},
			},
		},
		{
			Name:    "merge",
			Aliases: []string{"mg"},
			Action:  merge,
			Usage:   "Merge current working branch and cleanup everything",
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:    "develop",
					Aliases: []string{"dev"},
					Usage:   "Set merge as development to the default main branch",
				},
			},
		},
		{
			Name:    "submitreq",
			Aliases: []string{"sr"},
			Action:  submit,
			Usage:   "Create request to submit source back to Project",
		},
		{
			Name:    "close",
			Aliases: []string{"cl"},
			Action:  closeBranch,
			Usage:   "Close all work on this package branch (fall-back to the main branch, delete current)",
		},
		{
			Name:    "info",
			Aliases: []string{"f"},
			Action:  packageInfo,
			Usage:   "Display general information about this package",
		},
	}

	if err := app.Run(os.Args); err != nil {
		wzlib_logger.GetCurrentLogger().Errorf("Error: %s", err.Error())
	}
}
