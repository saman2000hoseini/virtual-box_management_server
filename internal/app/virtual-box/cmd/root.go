package cmd

import (
	"CloudComputing/virtual-box/internal/app/virtual-box/cmd/server"
	"github.com/spf13/cobra"

	"CloudComputing/virtual-box/internal/app/virtual-box/config"
)

// NewRootCommand creates a new virtual-box root command.
func NewRootCommand() *cobra.Command {
	var root = &cobra.Command{
		Use: "virtual-box",
	}

	cfg := config.Init()

	server.Register(root, cfg)

	return root
}
