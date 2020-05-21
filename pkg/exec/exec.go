package exec

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/den-is/ktempl/pkg/logging"
)

func ExecCommand(command string) {

	cmdSlice := strings.Fields(command)

	// TODO: central validation
	_, pathErr := exec.LookPath(cmdSlice[0])
	if pathErr != nil {
		logging.LogWithFields(
			logging.Fields{
				"component": "exec",
			}, "error", fmt.Sprintf("Was not able to find %q in $PATH. ", command), pathErr)
		os.Exit(1)
	}

	cmd := exec.Command(cmdSlice[0], cmdSlice[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		logging.LogWithFields(
			logging.Fields{
				"component": "exec",
			}, "error", "cmd.Start() failed: ", err)
	}

	if err := cmd.Wait(); err != nil {
		logging.LogWithFields(
			logging.Fields{
				"component": "exec",
			}, "error", "cmd.Start() ended with error: ", err)
	} else {
		fmt.Println("ExecCommand finished successfully")
		logging.LogWithFields(
			logging.Fields{
				"component": "exec",
			}, "info", "ExecCommand finished successfully")
	}

}
