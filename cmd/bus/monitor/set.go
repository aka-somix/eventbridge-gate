package monitor

import (
	"fmt"

	"github.com/spf13/cobra"
)

// setCmd
var setCmd = &cobra.Command{
	Use:   "set <busname>",
	Short: "Sets the monitor for a specific bus",
	Long:  `This command will create the resources necessary to start monitoring the Eventbridge Bus.`,
	Args:  cobra.ExactArgs(1), // Ensures exactly one argument is passed
	Run: func(cmd *cobra.Command, args []string) {
		busName := args[0]
		fmt.Printf("Monitor set for bus provided: %s\n", busName)
	},
}
