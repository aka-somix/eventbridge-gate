package profile

import (
	"github.com/spf13/cobra"
)

// profilesCmd represents the "profiles" subcommand
var Cmd = &cobra.Command{
	Use:   "profile",
	Short: "Manage the active AWS CLI profile",
	Long:  `This command allows you to set, view, or clear the active AWS CLI profile.`,
}

func init() {
	Cmd.AddCommand(setProfileCmd)
}
