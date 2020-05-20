package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	version = "dev"
	commit  = "dev"
	date    = "n/a"

	short_version bool

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
	versionCmd.Flags().BoolVarP(&short_version, "short", "s", false, "Return only version string")
}

func getVersion() string {

	if short_version {
		return version
	}
	return fmt.Sprintf("version: %s\ncommit: %s\ndate: %s\nlicense: %s\noriginal author: %s", version, commit, date, license, author)

}
