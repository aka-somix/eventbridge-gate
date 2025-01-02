package bus

import (
	"encoding/json"
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"
)

// listCmd represents the "list" subcommand
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists AWS EventBridge event buses",
	Long:  `This command is used to list the AWS EventBridge event buses.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Run the AWS CLI command to list event buses
		out, err := exec.Command("aws", "events", "list-event-buses", "--output", "json").Output()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		// Parse the JSON output
		var result struct {
			EventBuses []EventBus `json:"EventBuses"`
		}
		if err := json.Unmarshal(out, &result); err != nil {
			fmt.Printf("Error parsing AWS CLI output: %v\n", err)
			return
		}

		// Display the event buses
		if len(result.EventBuses) == 0 {
			fmt.Println("No EventBridge event buses found.")
			return
		}

		fmt.Println("Available EventBridge event buses:")
		for _, bus := range result.EventBuses {
			fmt.Printf(" - %s\n", bus.Name)
		}
	},
}
