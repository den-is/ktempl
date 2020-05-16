package cmd

import (
	"fmt"
	"os"

	"github.com/den-is/ktempl/pkg/logging"
)

func CheckFileExists(fname string) error {

	if _, err := os.Stat(fname); err == nil {

		return nil

	} else if os.IsNotExist(err) {

		logging.LogWithFields(
			logging.Fields{
				"component": "validation",
			}, "error", "File ", err)
		return err

	} else {

		logging.LogWithFields(
			logging.Fields{
				"component": "validation",
			}, "error", fmt.Sprintf("Issue with %q", fname), err)
		return err

	}

}
