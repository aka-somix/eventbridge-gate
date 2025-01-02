package profile

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

// Set the active profile using a wizard
var setProfileCmd = &cobra.Command{
	Use:   "set",
	Short: "Set the active AWS CLI profile",
	Long:  `Interactively select an AWS CLI profile from the list of available profiles to set as the active profile.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Fetch available profiles
		out, err := exec.Command("aws", "configure", "list-profiles").Output()
		if err != nil {
			fmt.Printf("Error fetching profiles: %v\n", err)
			return
		}

		// Parse profiles into a slice
		profiles := strings.Split(string(bytes.TrimSpace(out)), "\n")
		if len(profiles) == 0 {
			fmt.Println("No AWS profiles found.")
			return
		}

		// Use promptui for interactive selection
		prompt := promptui.Select{
			Label: "Select an AWS Profile",
			Items: profiles,
		}

		_, selectedProfile, err := prompt.Run()
		if err != nil {
			fmt.Printf("Error during selection: %v\n", err)
			return
		}

		// Set the selected profile as the active profile
		profileStore.SetActiveProfile(selectedProfile)
		fmt.Printf("Active profile set to '%s'\n", selectedProfile)
	},
}
