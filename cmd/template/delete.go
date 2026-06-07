package template

import (
	"fmt"

	"github.com/aaron70/decoy-cli/cli"
	"github.com/spf13/cobra"
)

func createDeleteCommand(cli *cli.CLI) *cobra.Command {
	command := &cobra.Command{
		Use: "delete <name>",
		Short: "Deletes the given template.",
		Example: `# Delete a template
decoy template delete "greet"`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]

			entity, err := cli.TemplateSvc.Delete(name)
			if err != nil {
				return err
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Template %q has been deleted:\n%s\n", entity.Name, entity.Tmpl)
			return nil
		},
	}

	return command
}
