package monitor

import (
	"fmt"

	"github.com/spf13/cobra"
)

// setCmd
var unsetCmd = &cobra.Command{
	Use:   "unset <busname>",
	Short: "Removes the monitor for a specific bus",
	Long:  `This command will destroy the resources necessary to monitor the Eventbridge Bus.`,
	Args:  cobra.ExactArgs(1), // Ensures exactly one argument is passed
	Run: func(cmd *cobra.Command, args []string) {
		// Extract bus name from args
		busName := args[0]
		// Create the monitor
		monitorService.Destroy(busName)
		// TODO scirone: beautify
		fmt.Printf("Monitor unset for bus provided: %s\n", busName)
	},
}
