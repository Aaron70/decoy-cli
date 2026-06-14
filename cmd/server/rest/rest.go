package rest

import (
	"github.com/aaron70/decoy-cli/internal/services"
	"github.com/spf13/cobra"
)

func CreateRestCommand(decoy *services.Decoy) *cobra.Command {
	cmd := &cobra.Command{
		Use: "rest",
		Short: "Manages servers of type REST",
	}

	cmd.AddCommand(CreateStartCommand(decoy))

	return cmd
}
