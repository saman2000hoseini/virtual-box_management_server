package cmd

import (
	"github.com/saman2000hoseini/virtual-box_management_server/internal/app/virtual-box/cmd/server"
	"github.com/saman2000hoseini/virtual-box_management_server/internal/app/virtual-box/config"
	"github.com/spf13/cobra"
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
