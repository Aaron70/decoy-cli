package template

import (
	"github.com/aaron70/decoy-cli/cli"
	"github.com/spf13/cobra"
)

func CreateTemplateCommand(cli *cli.CLI) *cobra.Command {
	command := &cobra.Command{
		Use:   "template",
		Short: "Groups the templates commands",
	}

	command.AddCommand(createStoreCommand(cli))
	command.AddCommand(createGetCommand(cli))
	command.AddCommand(createDeleteCommand(cli))
	command.AddCommand(createParseCommand(cli))

	return command
}
