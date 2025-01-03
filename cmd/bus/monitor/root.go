package monitor

import (
	"github.com/spf13/cobra"
)

// busCmd represents the "bus" subcommand
var Cmd = &cobra.Command{
	Use:   "monitor",
	Short: "Manages Monitors on event buses",
	Long:  `This command is used to manage monitors applied to AWS EventBridge event buses.`,
}

func init() {
	Cmd.AddCommand(monitorListCmd, setCmd, unsetCmd, tailCmd)
}
