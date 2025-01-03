package bus

import (
	"fmt"

	"github.com/spf13/cobra"
)

// listCmd represents the "list" subcommand
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists AWS EventBridge event buses",
	Long:  `This command is used to list the AWS EventBridge event buses.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Retrieve Event Buses
		eventBuses, err := ebService.List()
		
		if err != nil {
			// TODO scirone: manage errors
			return
		}

		// TODO scirone: beautify

		fmt.Println("Available EventBridge event buses:")
		for _, bus := range eventBuses {
			fmt.Printf(" - %s\n", bus.Name)
		}
	},
}
