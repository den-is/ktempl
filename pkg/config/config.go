package config

import (
	"fmt"
	"os"
	"time"

	"github.com/den-is/ktempl/pkg/logging"
	"github.com/spf13/viper"
)

type Config struct {
	Kubeconfig  string
	Timeout     string
	Retries     int
	Selector    string
	Namespace   string
	Template    string
	Output      string
	Permissions uint32
	Interval    time.Duration
	Exec        string
	Values      map[string]string
	Daemon      bool
	Pods        bool
	Log         logging.LoggingConfig
}

var (
	CurrentConfig Config
)

func Init() {

	viper.SetDefault("Permissions", 0644)
	viper.SetDefault("log.file", "")
	viper.SetDefault("log.level", "info")
	viper.SetDefault("timeout", "10s")
	viper.SetDefault("retries", "3")

	cfgFile := viper.GetString("config")

	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
	}

	viper.AddConfigPath(".")
	viper.AddConfigPath("/etc/ktempl")

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("No config file provided")
	} else {
		fmt.Println("Using config file:", cfgFile)
	}

	if err := viper.Unmarshal(&CurrentConfig); err != nil {
		fmt.Printf("Unable to decode into struct, %v", err)
		os.Exit(1)
	}

	logging.LoggerSetup(&CurrentConfig.Log)
	logging.LogWithFields(logging.Fields{
		"component": "config",
	}, "info", "Successfully initialized configuration")

}
