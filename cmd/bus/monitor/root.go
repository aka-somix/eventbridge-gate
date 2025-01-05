package monitor

import (
	"github.com/aka-somix/eventbridge-gate/pkg/services"
	"github.com/spf13/cobra"
)

// Global Monitor Service init
var monitorService, _ = services.NewMonitorService()

// busCmd represents the "bus" subcommand
var Cmd = &cobra.Command{
	Use:   "monitor",
	Short: "Manages Monitors on event buses",
	Long:  `This command is used to manage monitors applied to AWS EventBridge event buses.`,
}

func init() {
	Cmd.AddCommand(monitorListCmd, setCmd, unsetCmd, tailCmd)
}
