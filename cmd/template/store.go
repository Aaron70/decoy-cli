package template

import (
	"errors"
	"fmt"
	"os"

	"github.com/aaron70/decoy-cli/cli"
	errs "github.com/aaron70/goaty/errors"
	"github.com/spf13/cobra"
)

func createStoreCommand(cli *cli.CLI) *cobra.Command {
	var (
		name, tmpl, file string
		err              error
	)
	command := &cobra.Command{
		Use:  "store <name>",
		Short: "Upserts the contents of a template.",
		Long: "Upserts the contents of a template. You can pass the templates content from stdin, a file or trough the template flag.",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name = args[0]
			if cmd.Flags().Changed("file") {
				bytes, err := os.ReadFile(file)
				if err != nil {
					return err
				}
				tmpl = string(bytes)
			} else if !cmd.Flags().Changed("template") {
				tmpl, err = cli.ReadStringFrom(cmd.InOrStdin())
				if err != nil {
					return fmt.Errorf("Couldn't read the template from stdin: %w", err)
				}
			}

			oldTmpl, err := cli.TemplateSvc.Get(name)
			if err != nil {
				if !errors.Is(err, errs.ErrNotFound){
					return err
				}
			} else {
				fmt.Printf("Updating template %q:\n%s\n", name, oldTmpl)
				entity, err := cli.TemplateSvc.Update(name, tmpl)
				if err != nil {
					return err
				}
				tmpl = entity.Tmpl
				fmt.Printf("to:\n%s\n", tmpl)
				return nil
			}

			entity, err := cli.TemplateSvc.Save(name, tmpl)
			fmt.Printf("Saving new template %q:\n%s\n", name, entity.Tmpl)

			return err
		},
	}

	command.Flags().StringVarP(&tmpl, "template", "t", "", "The content of the template")
	command.Flags().StringVarP(&file, "file", "f", "", "The path of the file with the template contents")
	command.MarkFlagsMutuallyExclusive("template", "file")

	return command
}
