package monitor

import (
	"fmt"

	"github.com/spf13/cobra"
)

// listCmd represents the "list" subcommand
var monitorListCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists AWS EventBridge event buses",
	Long:  `This command is used to list the AWS EventBridge event buses.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Retrieve list of monitors from account
		var monitors = monitorService.List();

		// TODO scirone: beautify

		// Print out monitors found
		fmt.Printf("Monitors available: %s\n", monitors)
	},
}
