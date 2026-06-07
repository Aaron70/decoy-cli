package runner

import (
	"fmt"

	"github.com/aaron70/decoy-cli/cli"
	"github.com/spf13/cobra"
)

func createGetCommand(cli *cli.CLI) *cobra.Command {
	command := &cobra.Command{
		Use:   "get [<name>]",
		Args:  cobra.MaximumNArgs(1),
		Short: "List all the runners or get the details of the given runner.",
		Example: `# List all runners
decoy runner get

# Get the config of the given runner
decoy runner get "echo"
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 0 {
				entity, err := cli.RunnerSvc.Get(args[0])
				if err != nil {
					return err
				}
				fmt.Fprintf(cmd.OutOrStdout(), "%s\n", string(entity.Config))
			} else {
				entities, err := cli.RunnerSvc.GetAll()
				if err != nil {
					return err
				}
				for _, entity := range entities {
					fmt.Fprintf(cmd.OutOrStdout(), "%s\n", entity.Name)
				}
			}
			return nil
		},
	}
	return command
}
