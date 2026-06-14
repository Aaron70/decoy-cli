package rest

import (
	"os"

	"github.com/aaron70/decoy-cli/internal/services"
	server "github.com/aaron70/decoy-cli/server/rest"
	"github.com/spf13/cobra"
)

func CreateStartCommand(decoy *services.Decoy) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "Starts a new server to listen for requests",
		RunE: func(cmd *cobra.Command, args []string) error {
			specFile, err := os.OpenFile("./spec.yaml", os.O_RDONLY, os.ModePerm)
			if err != nil {
				return err
			}
			return server.Start(cmd.OutOrStdout(), decoy, 8080, specFile)
		},
	}
	return cmd
}
