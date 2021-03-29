package server

import (
	"CloudComputing/virtual-box/internal/app/virtual-box/config"
	"github.com/spf13/cobra"
)

func main(cfg config.Config) {

}

// Register registers server command for virtual-box binary.
func Register(root *cobra.Command, cfg config.Config) {
	root.AddCommand(
		&cobra.Command{
			Use:   "server",
			Short: "Run virtual-box server component",
			Run: func(cmd *cobra.Command, args []string) {
				main(cfg)
			},
		},
	)
}
