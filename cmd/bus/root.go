package bus

import (
	"github.com/spf13/cobra"
)

type EventBus struct {
	Name string `json:"Name"`
}

// busCmd represents the "bus" subcommand
var busCmd = &cobra.Command{
	Use:   "bus",
	Short: "Manages AWS EventBridge event buses",
	Long:  `This command is used to manage the AWS EventBridge event buses.`,
}

func init() {
	busCmd.AddCommand(listCmd)
}
