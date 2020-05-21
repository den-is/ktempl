package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	version = "dev"
	commit  = "dev"
	date    = "n/a"

	shortVersion bool

	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Get version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(getVersion())
		},
	}
)

const (
	license = "MPL 2.0"
	author  = "Denis Iskandarov denis@cloudboom.io"
)

func init() {
	rootCmd.AddCommand(versionCmd)
	versionCmd.Flags().BoolVarP(&shortVersion, "short", "s", false, "Return only version string")
}

func getVersion() string {

	if shortVersion {
		return version
	}
	return fmt.Sprintf("version: %s\ncommit: %s\ndate: %s\nlicense: %s\noriginal author: %s", version, commit, date, license, author)

}
