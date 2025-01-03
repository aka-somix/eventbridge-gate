package bus

import (
	"github.com/aka-somix/aws-events-gate/cmd/bus/monitor"
	"github.com/spf13/cobra"
)

type EventBus struct {
	Name string `json:"Name"`
}

// busCmd represents the "bus" subcommand
var Cmd = &cobra.Command{
	Use:   "bus",
	Short: "Manages AWS EventBridge event buses",
	Long:  `This command is used to manage the AWS EventBridge event buses.`,
}

func init() {
	Cmd.AddCommand(listCmd, monitor.Cmd)
}
