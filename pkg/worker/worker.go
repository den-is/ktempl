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
	selector := render.StringifyStringsMap(viper.GetStringMapString("selector"))
	output := viper.GetString("output")
	usePods := viper.GetBool("pods")

	template := viper.GetString("template")

	templateData := render.TemplData{}

	userProvidedValues := viper.GetStringMap("values")
	templateData.Values = &userProvidedValues

	for {

		conn, err := kubernetes.Connect(&kubeconfig)
		if err != nil {
			logging.LogWithFields(
				logging.Fields{
					"component": "runner",
				}, "error", "Connection to kubernetes has failed.", err)
			os.Exit(1)
		}

		nodes := kubernetes.GetHostList(conn, &namespace, &selector, &usePods)

		templateData.Nodes = nodes

		if err := render.ProduceOutput(template, &templateData, output); err == nil {
			if viper.GetString("exec") != "" {
				fmt.Println("Going to execute provided command")
				exec.ExecCommand(viper.GetString("exec"))
			}
		}

		// check if ktempl runs as a service or just once
		if viper.GetBool("daemon") {

			time.Sleep(viper.GetDuration("interval"))

		} else {
			os.Exit(0)
		}

	}

}
