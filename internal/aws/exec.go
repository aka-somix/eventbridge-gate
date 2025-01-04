package aws

import (
	"os/exec"

	"github.com/aka-somix/aws-events-gate/internal/store"
)

func AwsCommand(name string, arg ...string) *exec.Cmd {
	// Retrieve Profile
	profile, _ := store.GetProfileStore().GetActiveProfile()

	var cmdArgs []string
	cmdArgs = append(cmdArgs, arg...)

	// Add --profile if profile is not nil
	if profile != "" {
		// fmt.Printf("Executing request with profile: %v\n", profile)
		cmdArgs = append(cmdArgs, "--profile", profile)
	} else {
		// fmt.Println("Executing request with DEFAULT profile")
	}
	// Execute command
	return exec.Command(name, cmdArgs...)
}