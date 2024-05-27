package cmd

import (
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/qmstar0/shutdown"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"simple-gateway/cmd/middlewate"
	"simple-gateway/cmd/router"
	"simple-gateway/config"
	"simple-gateway/proxy"

	"net/http"
	"os"
)

const defaultConfigPath = "proxy.toml"

var (
	configPath string
	debug      bool
)

var (
	cfg      config.Config
	proxyMap = make(map[string]http.Handler)
)

var rootCmd = &cobra.Command{
	Use:   "proxy",
	Short: "Forward the request to the corresponding service",
	Long:  "Forward the request to the corresponding service, which comes with a static resource service",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {

		if debug {
			log.SetLevel(log.DebugLevel)
		}

		if len(cfg.Proxys) <= 0 {
			log.Error("没有配置任何反向代理")
			shutdown.Exit(1)
		}
		for _, p := range cfg.Proxys {
			log.Infof("Source:%s › Target:%s", p.Source, p.Target)
			host, handler, err := proxy.NewReverseProxy(p.Source, p.Target)
			if err != nil {
				log.Error(err)
				shutdown.Exit(1)
			}
			if err = SetUpHttpHandler(host, handler); err != nil {
				log.Error(err)
				shutdown.Exit(1)
			}
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		logger := log.WithPrefix("Reverse proxy server")
		router.SetUpRouter(proxyMap)
		middlewate.SetUpMiddlewate(proxyMap)

		for _, p := range cfg.Port {
			go func(port int) {
				addr := fmt.Sprintf("0.0.0.0:%d", port)
				logger.Infof("Start listening %s", addr)
				logger.Info(http.ListenAndServe(addr, nil))
			}(p)
		}
	},
}

func SetUpHttpHandler(host string, handler http.Handler) error {
	if _, ok := proxyMap[host]; ok {
		return fmt.Errorf("重复配置同一地址")
	}
	proxyMap[host] = handler
	return nil
}

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.PersistentFlags().StringVar(&configPath, "config", defaultConfigPath, "Set configuration file path")
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "debug mode outputs more information")
	cobra.OnInitialize(InitConfig)
}

func InitConfig() {
	viper.SetConfigFile(configPath)
	if err := viper.ReadInConfig(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
		return
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
		return
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%s", err)
		shutdown.Exit(1)
	}
	shutdown.WaitCtrlC()
}
