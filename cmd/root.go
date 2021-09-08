package cmd

import (
	"github.com/spf13/cobra"
	"github.com/wmw9/rekoda/internal/config"
)

//var CfgFile string

func NewRootCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "rekoda",
		Short: "📼  " + APP + "v" + VERSION + "- Automatic Twitch Recorder",
		Long:  "📼  " + APP + " v" + VERSION + " - Automatic Twitch Recorder. Homepage: https://github.com/wmw9/rekoda",
	}
}

var rootCmd = NewRootCmd()

type Config struct {
	Title    string `toml:"title"`
	Version  int    `toml:"version"`
	Channels struct {
		Reckful struct {
			ID      int64  `toml:"id"`
			Enabled bool   `toml:"enabled"`
			Quality string `toml:"quality"`
		} `toml:"reckful"`
	} `toml:"channels"`
	Tokens struct {
		Telegram struct {
			LogChannelEnabled bool   `toml:"log_channel_enabled"`
			LogChannelID      int    `toml:"log_channel_id"`
			BotToken          string `toml:"bot_token"`
		} `toml:"telegram"`
		Twitch struct {
			ClientID string `toml:"client_id"`
		} `toml:"twitch"`
	} `toml:"tokens"`
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	//rootCmd := NewRootCmd()
	//	cobra.OnInitialize(config.InitConfig)

	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.PersistentFlags().StringVarP(&config.FlagLogLevel, "verbose", "v", "", "Set log level: trace, debug, info (default is 'info')")
	rootCmd.PersistentFlags().StringVarP(&config.FlagConfigFile, "config", "c", "", "Custom config file (default is $HOME/rekoda/rekoda.toml)")
	rootCmd.PersistentFlags().StringVarP(&config.FlagStreamsDir, "output", "o", "", "Custom stream directory to download (default is $HOME/rekoda/streams)")
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
