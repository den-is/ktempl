package worker

import (
	"fmt"
	"os"
	"time"

	"github.com/den-is/ktempl/pkg/exec"
	"github.com/den-is/ktempl/pkg/kubernetes"
	"github.com/den-is/ktempl/pkg/logging"
	"github.com/den-is/ktempl/pkg/render"
	"github.com/spf13/viper"
)

func Worker() {

	kubeconfig := viper.GetString("kubeconfig")
	namespace := viper.GetString("namespace")
	selector := viper.GetString("selector")
	output := viper.GetString("output")
	use_pods := viper.GetBool("pods")

	template_data := render.TemplData{}
	template := viper.GetString("template")

	// TODO: accept complex values for the right side of the key=value expression, rather than just string values
	// parse user provided values into map
	user_values := render.StringSliceToStringMap(viper.GetStringSlice("values"))
	template_data.Values = &user_values

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
			if viper.GetString("exec") != "" {
				fmt.Println("Going to execute provided command")
				exec.ExecCommand(viper.GetString("exec"))
			}
		}

		// check if ktempl runs as a service or just once
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
