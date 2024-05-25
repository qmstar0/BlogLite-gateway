package cmd

import (
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"simple-gateway/assets"
)

var assetsCmd = &cobra.Command{
	Use:   "assets",
	Short: "Use static file resources",
	Long:  "Use static file resources",
	PreRun: func(cmd *cobra.Command, args []string) {
		if err := viper.BindPFlags(cmd.Flags()); err != nil {
			log.Error(err)
			shutdown.Exit(1)
		}
		if err := viper.Unmarshal(&cfg); err != nil {
			log.Error(err)
			shutdown.Exit(1)
		}

		log.Infof("Source:%s › Assets Dir:%s", cfg.Assets.Source, cfg.Assets.Dir)
	},
	Run: func(cmd *cobra.Command, args []string) {
		host, handler, err := assets.NewAssetsFileServer(cfg.Assets)
		if err != nil {
			log.Error(err)
			shutdown.Exit(1)
		}
		if err = SetUpHttpHandler(host, handler); err != nil {
			log.Error(err)
			shutdown.Exit(1)
		}
		rootCmd.Run(cmd, args)
	},
}

func init() {
	assetsCmd.Flags().String("source", "", "Set the source address that requires a proxy")
	assetsCmd.Flags().String("dir", "", "Set static resource service dir")
	rootCmd.AddCommand(assetsCmd)
}
