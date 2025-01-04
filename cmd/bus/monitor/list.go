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
		var monitors, err = monitorService.List()
		if err != nil {
			fmt.Printf("Error while retrieving monitors list: %s\n", err)
			return
		}

		// Check if no monitors were found
		if len(monitors) == 0 {
			fmt.Println("No monitors found.")
			return
		}

		// TODO scirone: Beautify

		// Print out monitors in a bullet list format
		fmt.Println("Monitors available:")
		for _, monitor := range monitors {
			fmt.Printf("- %s\n", monitor)
		}
	},
}
