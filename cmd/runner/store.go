package runner

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
		name, config, file string
		err                      error
	)
	command := &cobra.Command{
		Use:   "store <name>",
		Short: "Upserts a runner.",
		Long:  "Upserts a runner. You can pass the config JSON from stdin, a file or trough the config flag.",
		Example: `# Store a runner config from stdin
echo 'echo "{{ .Template }}"' | decoy runner store "echo"

# Store a runner config from a file
decoy runner store "echo" -f /path/to/config

# Store an inline runner config
decoy runner store "echo" -c 'echo "{{ .Template }}"'
`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name = args[0]

			if cmd.Flags().Changed("file") {
				bytes, err := os.ReadFile(file)
				if err != nil {
					return err
				}
				config = string(bytes)
			} else if !cmd.Flags().Changed("config") {
				config, err = cli.ReadStringFrom(cmd.InOrStdin())
				if err != nil {
					return fmt.Errorf("Couldn't read the config from stdin: %w", err)
				}
			}

			oldRunner, err := cli.RunnerSvc.Get(name)
			if err != nil {
				if !errors.Is(err, errs.ErrNotFound) {
					return err
				}
			} else {
				fmt.Printf("Updating runner %q:\n%s\n", name, string(oldRunner.Config))
				entity, err := cli.RunnerSvc.Update(name, config)
				if err != nil {
					return err
				}
				config = string(entity.Config)
				fmt.Printf("to:\n%s\n", config)
				return nil
			}

			entity, err := cli.RunnerSvc.Save(name, config)
			fmt.Printf("Saving new runner %q:\n%s\n", name, string(entity.Config))

			return err
		},
	}

	command.Flags().StringVarP(&config, "config", "c", "", "The config JSON of the runner")
	command.Flags().StringVarP(&file, "file", "f", "", "The path of the file with the runner config")
	command.MarkFlagsMutuallyExclusive("config", "file")

	return command
}
