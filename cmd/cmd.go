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

	defaultDuration, _ := time.ParseDuration("15s")

	rootCmd.Flags().StringP("kubeconfig", "k", "", "path to kubeconfig (default is \"~/.kube/config\")")
	rootCmd.Flags().StringP("namespace", "n", "", "used with -p. kubernetes namespace indicator (default is \"\" or what is provided in kubeconfig. usually \"default\")")
	rootCmd.Flags().StringP("selector", "l", "", "kubernetes label selectors (default is \"\" a.k.a everything)")
	rootCmd.Flags().StringP("template", "t", "", "input template file. Required")
	rootCmd.Flags().StringP("output", "o", "", "output file path (default stdout)")
	rootCmd.Flags().DurationP("interval", "i", defaultDuration, "used with -d. interval between polls (default 15s)")
	rootCmd.Flags().StringP("config", "c", "", "config file location. Optional.")
	rootCmd.Flags().StringP("exec", "e", "", "execute command on success. Optional")
	rootCmd.Flags().String("log-file", "", "log file destination. default stdout. Optional")
	rootCmd.Flags().String("log-level", "", "log file destination. default stdout. Optional")
	rootCmd.Flags().StringToString("set", nil, "provide additional template values of form stringKey=StringValue (can be multiple comma separated values key1=val1,key2=val2)")
	rootCmd.Flags().BoolP("pods", "p", false, "query pods for their nodes")
	rootCmd.Flags().BoolP("daemon", "d", false, "run as daemon")

	_ = rootCmd.MarkFlagRequired("template")

	_ = viper.BindPFlag("kubeconfig", rootCmd.Flags().Lookup("kubeconfig"))
	_ = viper.BindPFlag("namespace", rootCmd.Flags().Lookup("namespace"))
	_ = viper.BindPFlag("selector", rootCmd.Flags().Lookup("selector"))
	_ = viper.BindPFlag("template", rootCmd.Flags().Lookup("template"))
	_ = viper.BindPFlag("output", rootCmd.Flags().Lookup("output"))
	_ = viper.BindPFlag("interval", rootCmd.Flags().Lookup("interval"))
	_ = viper.BindPFlag("config", rootCmd.Flags().Lookup("config"))
	_ = viper.BindPFlag("exec", rootCmd.Flags().Lookup("exec"))
	_ = viper.BindPFlag("log.file", rootCmd.Flags().Lookup("log-file"))
	_ = viper.BindPFlag("log.level", rootCmd.Flags().Lookup("log-level"))
	_ = viper.BindPFlag("values", rootCmd.Flags().Lookup("set"))
	_ = viper.BindPFlag("pods", rootCmd.Flags().Lookup("pods"))
	_ = viper.BindPFlag("daemon", rootCmd.Flags().Lookup("daemon"))

}
