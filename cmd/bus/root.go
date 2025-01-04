package bus

import (
	"github.com/aka-somix/aws-events-gate/cmd/bus/monitor"
	"github.com/aka-somix/aws-events-gate/pkg/services"
	"github.com/spf13/cobra"
)

// Global Monitor Service init
var ebService, _ = services.NewEventBusService()


// busCmd represents the "bus" subcommand
var Cmd = &cobra.Command{
	Use:   "bus",
	Short: "Manages AWS EventBridge event buses",
	Long:  `This command is used to manage the AWS EventBridge event buses.`,
}

func init() {
	Cmd.AddCommand(listCmd, monitor.Cmd)
}
