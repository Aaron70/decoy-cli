package runner

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/aaron70/decoy-cli/cli"
	"github.com/aaron70/decoy-cli/internal/services"
	"github.com/spf13/cobra"
)

func createRunCommand(cli *cli.CLI) *cobra.Command {
	var (
		_type services.RunnerType
		tmpl, file, data, config string
		stringValues                    map[string]string
		n, g                               int
		err                             error
	)
	command := &cobra.Command{
		Use:  "run <type> [<runner>] [<template>]",
		Args: cobra.RangeArgs(1, 3),
		RunE: func(cmd *cobra.Command, args []string) error {
			_type = services.RunnerType(args[0])

			if len(args) >= 2 {
				entity, err := cli.RunnerSvc.Get(args[1])
				if err != nil {
					return err
				}
				config = entity.Config
			} else if !cmd.Flags().Changed("config") {
				fmt.Printf("reading from stdin\n")
				tmpl, err = cli.ReadStringFrom(cmd.InOrStdin())
				if err != nil {
					return fmt.Errorf("Couldn't read the runner's config from stdin: %w", err)
				}
			}

			if len(args) == 3 {
				entity, err := cli.TemplateSvc.Get(args[2])
				if err != nil {
					return err
				}
				tmpl = entity.Tmpl
			} else if cmd.Flags().Changed("file") {
				bytes, err := os.ReadFile(file)
				if err != nil {
					return err
				}
				tmpl = string(bytes)
			}

			var jsonData map[string]any
			err = json.Unmarshal([]byte(data), &jsonData)
			if err != nil {
				return fmt.Errorf("Couldn't deserialize the contents from the flag --data/-d: %w", err)
			}

			for key, value := range stringValues {
				jsonData[key] = value
			}

			return cli.RunnerSvc.Run(cmd.OutOrStdout(), _type, config, tmpl, jsonData, n, g)
		},
	}

	command.Flags().StringVarP(&tmpl, "template", "t", "", "The content of the template to parse")
	command.Flags().StringVarP(&file, "file", "f", "", "The path of the file with the template contents to parse")
	command.MarkFlagsMutuallyExclusive("template", "file")

	command.Flags().StringVarP(&data, "data", "d", "{}", "The JSON data to be used within the template")
	command.Flags().StringToStringVarP(&stringValues, "value", "v", map[string]string{}, "A set of pairs (key=value) to inject into the data")

	command.Flags().IntVarP(&n, "times", "n", 1, "The number of times that the runner will be executed")
	command.Flags().IntVarP(&g, "goroutines", "g", 1, "The number of concurrent goroutines executing the runner")

	command.Flags().StringVarP(&config, "config", "c", "", "The config's content of the runner to execute")

	return command
}
