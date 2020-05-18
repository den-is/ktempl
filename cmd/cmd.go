package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/den-is/ktempl/pkg/config"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	rootCmd = &cobra.Command{
		Use:   "ktempl",
		Short: "ktempl renders configs from templates with data received from a kubernetes cluster",
		Long: `ktempl queries Kubernetes cluster indicated in kubeconfig.
User included in kubeconfig should be allowed to at least list nodes in the cluster`,
		Example: `
# Select nodes with "disk=ssd" label selector
ktempl -l disk=ssd

# Selects nodes on which pods with label "app=myapp1" are running
ktempl -p -l app=myapp1 -t myconf_template.tpl -o myconf.yaml

ktempl -c config.yaml -t test.tpl  -o output.txt -l app=stagingapps --set title=XXX --log-level warn
`,
		Run: func(cmd *cobra.Command, args []string) {
			// delegate primary job execution control to dedicated function.
			// passing all arguments and command object.
			Gates(cmd, args)
		},
	}
)

func Start() {

	err := rootCmd.Execute()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}

func init() {

	// initializing viper
	cobra.OnInitialize(config.Init)

	default_duration, _ := time.ParseDuration("15s")

	rootCmd.PersistentFlags().StringP("kubeconfig", "k", "", "path to kubeconfig (default is \"~/.kube/config\")")
	rootCmd.PersistentFlags().StringP("namespace", "n", "", "used with -p. kubernetes namespace indicator (default is \"\" or what is provided in kubeconfig. usually \"default\")")
	rootCmd.PersistentFlags().StringP("selector", "l", "", "kubernetes label selectors (default is \"\" a.k.a everything)")
	rootCmd.PersistentFlags().StringP("template", "t", "", "input template file. Required")
	rootCmd.PersistentFlags().StringP("output", "o", "", "output file path (default stdout)")
	rootCmd.PersistentFlags().DurationP("interval", "i", default_duration, "used with -d. interval between polls (default 15s)")
	rootCmd.PersistentFlags().StringP("config", "c", "", "config file location. Optional.")
	rootCmd.PersistentFlags().StringP("exec", "e", "", "execute command on success. Optional")
	rootCmd.PersistentFlags().String("log-file", "", "log file destination. default stdout. Optional")
	rootCmd.PersistentFlags().String("log-level", "", "log file destination. default stdout. Optional")
	rootCmd.PersistentFlags().StringToString("set", nil, "provide additional template values of form stringKey=StringValue (can be multiple comma separated values key1=val1,key2=val2)")
	rootCmd.PersistentFlags().BoolP("pods", "p", false, "query pods for their nodes")
	rootCmd.PersistentFlags().BoolP("daemon", "d", false, "run as daemon")

	_ = rootCmd.MarkFlagRequired("template")

	_ = viper.BindPFlag("kubeconfig", rootCmd.PersistentFlags().Lookup("kubeconfig"))
	_ = viper.BindPFlag("namespace", rootCmd.PersistentFlags().Lookup("namespace"))
	_ = viper.BindPFlag("selector", rootCmd.PersistentFlags().Lookup("selector"))
	_ = viper.BindPFlag("template", rootCmd.PersistentFlags().Lookup("template"))
	_ = viper.BindPFlag("output", rootCmd.PersistentFlags().Lookup("output"))
	_ = viper.BindPFlag("interval", rootCmd.PersistentFlags().Lookup("interval"))
	_ = viper.BindPFlag("config", rootCmd.PersistentFlags().Lookup("config"))
	_ = viper.BindPFlag("exec", rootCmd.PersistentFlags().Lookup("exec"))
	_ = viper.BindPFlag("log.file", rootCmd.PersistentFlags().Lookup("log-file"))
	_ = viper.BindPFlag("log.level", rootCmd.PersistentFlags().Lookup("log-level"))
	_ = viper.BindPFlag("values", rootCmd.PersistentFlags().Lookup("set"))
	_ = viper.BindPFlag("pods", rootCmd.PersistentFlags().Lookup("pods"))
	_ = viper.BindPFlag("daemon", rootCmd.PersistentFlags().Lookup("daemon"))

}
