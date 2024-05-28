package cmd

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/qmstar0/nightsky-gateway/config"
	"github.com/qmstar0/nightsky-gateway/router"
	"github.com/qmstar0/nightsky-gateway/service"
	"github.com/qmstar0/shutdown"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"net/http"
	"os"
)

const defaultConfigPath = "config.toml"

var (
	configPath  string
	debug       bool
	cfg         config.Config
	proxyMap    = make(map[string]http.Handler)
	sourceStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#74D6FB")).Render
	targetStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#74FBB5")).Render
)

var rootCmd = &cobra.Command{
	Use:   "proxy",
	Short: "Forward the request to the corresponding service",
	Long:  "Forward the request to the corresponding service, which comes with a static resource service",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if debug || cfg.Debug {
			log.SetLevel(log.DebugLevel)
			log.Debug("debug mode")
		}

		for _, p := range cfg.Proxys {
			log.Infof("Source: %s  › Target: %s ", sourceStyle(p.Source), targetStyle(p.Target))
			handler, err := service.NewReverseProxy(p.Target)
			if err != nil {
				log.Error(err)
				shutdown.Exit(1)
			}
			if err := SetUpHttpHandler(p.Source, handler); err != nil {
				log.Error(err)
				shutdown.Exit(1)
			}
		}

		if cfg.Assets != nil {
			log.Infof("Source: %s  › AssetsDir: %s ", sourceStyle(cfg.Assets.Source), targetStyle(cfg.Assets.Dir))
			handler, err := service.NewAssetsFileServer(cfg.Assets)
			if err != nil {
				log.Error(err)
				shutdown.Exit(1)
			}
			if err := SetUpHttpHandler(cfg.Assets.Source, handler); err != nil {
				log.Error(err)
				shutdown.Exit(1)
			}
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		ssl := cfg.SSL != nil
		go router.ListenAndServe(proxyMap, ssl)

		if ssl {
			go router.ListenAndServeTLS(proxyMap, cfg.SSL.SSLCertFilePath, cfg.SSL.SSLKeyFilePath)
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
