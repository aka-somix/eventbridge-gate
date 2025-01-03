package monitor

import (
	"fmt"

	"github.com/spf13/cobra"
)

// setCmd
var tailCmd = &cobra.Command{
	Use:   "tail <busname>",
	Short: "Removes the monitor for a specific bus",
	Long:  `This command will tail the events emitted on the event bus`,
	Args:  cobra.ExactArgs(1), // Ensures exactly one argument is passed
	Run: func(cmd *cobra.Command, args []string) {
		// Extract bus name from args
		busName := args[0]
		// Create the monitor
		monitorService.Tail(busName)
		// TODO scirone: beautify
		fmt.Printf("No Available events to show for bus: %s\n", busName)
	},
}
