package root

import "github.com/spf13/cobra"

func CreateRootCmd() *cobra.Command {
	var rootCmd = &cobra.Command{
		Use:   "jc",
		Short: "jc short",
	}

	return rootCmd
}
