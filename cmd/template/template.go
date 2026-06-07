package template

import (
	"github.com/aaron70/decoy-cli/cli"
	"github.com/spf13/cobra"
)

func CreateTemplateCommand(cli *cli.CLI) *cobra.Command {
	command := &cobra.Command{
		Use:   "template",
		Short: "Groups the templates commands",
		Long: `Manages the templates commands, you can store, get, delete and parse templates.
The Template engine uses the Go's text/template library, so you can leverage the go tmpl syntax to structure and generate your data.`,
	}

	command.AddCommand(createStoreCommand(cli))
	command.AddCommand(createGetCommand(cli))
	command.AddCommand(createDeleteCommand(cli))
	command.AddCommand(CreateParseCommand(cli))

	return command
}
