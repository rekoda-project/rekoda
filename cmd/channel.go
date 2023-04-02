package cmd

import (
	"errors"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/rekoda-project/rekoda/internal/config"
)

var (
	errExactArgs    = errors.New("Too many channels, you may add only one channel at a time")
	errMinimumNArgs = errors.New("You need to specify at least one channel")
)

// channelCmd represents the channel command
func NewChannelCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "channel",
		Short: "Manage your channels: list, add, remove, disable or enable",
		Long:  `Manage your channels: list, add, remove, disable or enable for recording.`,
	}
}

var channelCmd = NewChannelCmd()

// NewAddCmd represents the channel command
func NewAddCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "add",
		Short: "Add channels to record",
		Long:  "Add channels to record",
		//Args:  cobra.MinimumNArgs(1),
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errMinimumNArgs
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			addChannels(args)
			return nil
		},
	}
}

var addCmd = NewAddCmd()

// NewRemoveCmd represents the channel command
func NewRemoveCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "remove",
		Short: "Remove channels from config",
		Long:  "Remove channels from config",
		//Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			removeChannels(args)
			return nil
		},
	}
}

var removeCmd = NewRemoveCmd()

// NewListCmd represents the channel command
func NewListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all channels",
		Long:  "List all channels listed in config file",
		//Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			listChannels()
			return nil
		},
	}
}

var listCmd = NewListCmd()

// NewDisableCmd represents the channel command
func NewDisableCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "disable",
		Short: "Disables channels from recording",
		Long:  "Disables channels from recording",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errMinimumNArgs
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			disableChannels(args)
			return nil
		},
	}
}

var disableCmd = NewDisableCmd()

// NewEnableCmd represents the channel command
func NewEnableCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "enable",
		Short: "Enables channels for recording",
		Long:  "Enables channels from recording",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errMinimumNArgs
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			enableChannels(args)
			return nil
		},
	}
}

var enableCmd = NewEnableCmd()

func init() {
	rootCmd.AddCommand(channelCmd)
	channelCmd.AddCommand(addCmd)
	channelCmd.AddCommand(removeCmd)
	channelCmd.AddCommand(listCmd)
	channelCmd.AddCommand(disableCmd)
	channelCmd.AddCommand(enableCmd)
}

func addChannels(list []string) {
	c := config.InitConfig()

	for _, v := range list {
		// NYI: Check if channel not banned/valid

		// Check if channel already exists in config
		if c.IsChannelInConfig(v) {
			log.WithField("general", "CLI").Warnf("Skip! '%v' is already added in config", v)
			continue
		}
		log.WithField("general", "CLI").Tracef("Trying to add '%v' in config", v)
		channel := config.Channels{
			User:    v,
			ID:      123,
			Enabled: true,
			Quality: "best",
		}
		c.Channels = append(c.Channels, channel)
		log.WithField("general", "CLI").Infof("Channel added '%v' in config file", v)
	}
	if err := c.Save(); err != nil {
		panic(err)
	}
}

func removeChannels(list []string) {
	c := config.InitConfig()

	for _, v := range list {
		c.Channels = removeSliceByName(c.Channels, v)
		log.WithField("general", "CLI").Infof("Channel '%v' removed from config file", v)
	}
	if err := c.Save(); err != nil {
		panic(err)
	}
}

func disableChannels(list []string) {
	c := config.InitConfig()

	for _, v := range list {
		for i, u := range c.Channels {
			if u.User == v {
				log.WithField("general", "CLI").Infof("Found! %v trying to disable it", u.User)
				c.Channels[i].Enabled = false
			}
		}
	}
	if err := c.Save(); err != nil {
		panic(err)
	}
}

func enableChannels(list []string) {
	c := config.InitConfig()

	for _, v := range list {
		for i, u := range c.Channels {
			if u.User == v {
				log.WithField("general", "CLI").Infof("Found! '%v' trying to enable it", u.User)
				c.Channels[i].Enabled = true
			}
		}
	}
	if err := c.Save(); err != nil {
		panic(err)
	}
}

func listChannels() {
	c := config.InitConfig()

	ctxLog := log.WithField("general", "LIST")
	total, rec := c.ChannelList()
	ctxLog.Info(total)
	ctxLog.Info(rec)
}

func removeSliceByName(s []config.Channels, r string) []config.Channels {
	for i, v := range s {
		if v.User == r {
			return append(s[:i], s[i+1:]...)
		}
	}
	return s
}
