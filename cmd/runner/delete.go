package runner

import (
	"fmt"

	"github.com/aaron70/decoy-cli/cli"
	"github.com/spf13/cobra"
)

func createDeleteCommand(cli *cli.CLI) *cobra.Command {
	command := &cobra.Command{
		Use:  "delete <name>",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]

			entity, err := cli.RunnerSvc.Delete(name)
			if err != nil {
				return err
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Runner %q has been deleted:\n%s\n", entity.Name, string(entity.Config))
			return nil
		},
	}

	return command
}
