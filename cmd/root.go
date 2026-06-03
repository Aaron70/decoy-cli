package cmd

import (
	"os"
	"path/filepath"

	"github.com/aaron70/decoy-cli/cli"
	"github.com/aaron70/decoy-cli/cmd/runner"
	"github.com/aaron70/decoy-cli/cmd/template"
	"github.com/spf13/cobra"
)

func createRootCommand(cli *cli.CLI) *cobra.Command {
	command := &cobra.Command{
		Use:   "decoy",
		Short: "A CLI tool to generate mock data",
		Long:  "Decoy is a command-line utility for creating and ingesting mock data trough templates and runners. Templates support dynamic data injection using the Go Template Engine, and runners let you ingest the generated data into your application.",
	}

	command.AddCommand(template.CreateTemplateCommand(cli))
	command.AddCommand(template.CreateParseCommand(cli))
	command.AddCommand(runner.CreateRunnerCommand(cli))
	command.AddCommand(runner.CreateRunCommand(cli))

	return command
}

func configDir() string {
		var basePath string
    if dir := os.Getenv("XDG_CONFIG_HOME"); dir != "" {
        basePath = dir
    }
    home, _ := os.UserHomeDir()
		basePath = home
    return filepath.Join(basePath, ".config", "decoy")
}

func Execute() error {
	cli, err := cli.NewCLI(configDir())
	if err != nil {
		return err
	}
	err = createRootCommand(cli).Execute()
	if err != nil {
		return err
	}
	return nil
}
