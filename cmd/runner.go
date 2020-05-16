package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/den-is/ktempl/pkg/exec"
	"github.com/den-is/ktempl/pkg/kubernetes"
	"github.com/den-is/ktempl/pkg/logging"
	"github.com/den-is/ktempl/pkg/render"
	"github.com/den-is/ktempl/pkg/validation"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Starts primary job
func StartJob(cmd *cobra.Command, args []string) {

	kubeconfig := viper.GetString("kubeconfig")
	namespace := viper.GetString("namespace")
	selector := viper.GetString("selector")
	template := viper.GetString("template")
	output := viper.GetString("output")
	use_pods := viper.GetBool("pods")

	template_data := render.TplData{}

	// TODO: accept complex values for the right side of the key=value expression, rather than just string values
	// parse user provided values into map
	user_values := render.StringSliceToStringMap(viper.GetStringSlice("set"))
	template_data.Values = &user_values

	// TODO: add central place for validation logic

	// Check if template file exists
	if err := validation.CheckFileExists(template); err != nil {
		logging.LogWithFields(
			logging.Fields{
				"component": "runner",
			}, "error", "Template file at given path does not exist.", err)
		os.Exit(1)
	}

	for {

		conn, err := kubernetes.Connect(&kubeconfig)
		if err != nil {
			logging.LogWithFields(
				logging.Fields{
					"component": "runner",
				}, "error", "Connection to kubernetes has failed.", err)
			os.Exit(1)
		}

		nodes := kubernetes.GetHostList(conn, &namespace, &selector, &use_pods)

		template_data.Nodes = nodes

		if err := render.RenderOutput(template, &template_data, output); err == nil {
			fmt.Println("Going to execute command")
			if viper.GetString("exec") != "" {
				exec.ExecCommand(viper.GetString("exec"))
			}
		}

		if viper.GetBool("daemon") {
			if interval_duration, err := time.ParseDuration(viper.GetString("interval")); err == nil {
				time.Sleep(interval_duration)
			} else {
				logging.LogWithFields(
					logging.Fields{
						"component": "runner",
					}, "error", "Was not able to parse user provided duration ", viper.GetString("interval"))
				os.Exit(1)
			}
		} else {
			os.Exit(0)
		}

	}

}
