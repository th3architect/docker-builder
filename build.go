package main

import (
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/sylphon/build-runner"
	"github.com/sylphon/build-runner/unit-config"

	"github.com/codegangsta/cli"
)

func build(c *cli.Context) {
	builderfile := c.Args().First()
	if builderfile == "" {
		builderfile = "Bobfile"
	}

	unitConfig, err := unitconfig.ReadFromFile("./"+builderfile, unitconfig.TOML)
	if err != nil {
		if c.Bool("force") {
			err = forceBuild()
		}
	}
	if err != nil {
		exitErr(1, "unable to parse unit config", err)
	}

	unitConfig.SkipPush(c.Bool("skip-push"))

	if err := runner.RunBuild(unitConfig, os.Getenv("PWD")); err != nil {
		exitErr(1, "unable to build", err)
	}

}

func forceBuild() error {
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}

	basename := path.Base(pwd)
	args := []string{"docker", "build", "-t", basename, "."}

	cmd := exec.Command("docker")
	cmd.Args = args
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin

	Logger.Info("running command --> " + strings.Join(args, " "))
	return cmd.Run()
}
