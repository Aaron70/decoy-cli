package runner

import (
	"github.com/aaron70/decoy-cli/cli"
	"github.com/spf13/cobra"
)

func CreateRunnerCommand(cli *cli.CLI) *cobra.Command {
	command := &cobra.Command{
		Use:   "runner",
		Short: "Groups the runner commands",
	}

	command.AddCommand(createStoreCommand(cli))
	command.AddCommand(createGetCommand(cli))
	command.AddCommand(createDeleteCommand(cli))
	command.AddCommand(CreateRunCommand(cli))

	return command
}
