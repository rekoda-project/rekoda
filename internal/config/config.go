package config

import (
	"fmt"
	"os"

	"github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	conf "github.com/rekoda-project/rekoda/pkg/config/toml"
)

const ConfigFile = "rekoda.toml"

var (
	sep            = string(os.PathSeparator)
	homeDir        = ""
	configDir      = ""
	streamsDir     = ""
	EnvConfigFile  = os.Getenv("REKODA_CONF_FILE")
	EnvConfigDir   = os.Getenv("REKODA_CONF_DIR")
	EnvStreamsDir  = os.Getenv("REKODA_STREAMS_DIR")
	FlagConfigFile = ""
	FlagConfigDir  = ""
	FlagStreamsDir = ""
	EnvLogLevel    = os.Getenv("REKODA_LOG_LEVEL") // Options: trace, debug, info (default)
	FlagLogLevel   = ""
)

type Config struct {
	Title      string     `toml:"title"`
	Version    int        `toml:"version"`
	ConfigDir  string     `toml:"config_dir"`
	ConfigFile string     `toml:"config_file"`
	StreamsDir string     `toml:"streams_dir"`
	Channels   []Channels `toml:"channels"`
}

type Channels struct {
	Enabled bool   `toml:"enabled"`
	User    string `toml:"user"`
	ID      int64  `toml:"id"`
	Quality string `toml:"quality"`
}

var ConfigStruct Config

// MakeDefaultConfStruct initializes default conf struct
func (c *Config) MakeDefaultConfStruct() {
	c.Title = "Rekoda configuration file"
	c.Version = 1
	c.ConfigDir = configDir
	c.StreamsDir = streamsDir
	c.Channels = []Channels{
		{
			Enabled: false,
			User:    "test",
			ID:      999999999999,
			Quality: "best",
		},
	}
}

// initConfig reads in config file and ENV variables if set, otherwise initialize default conf
func InitConfig() *Config {
	ctxLog := log.WithField("general", "INIT")

	c := &Config{}
	c.AutomaticEnv(ctxLog)

	// At this step config file should be ready, we open it
	ctxLog.Debugf("Trying to open config file: %s", c.ConfigFile)
	c.StreamsDir = c.ConfigDir + sep + "streams"
	ctxLog.Tracef("Config streams dir: '%v'", c.StreamsDir)
	c.Channels = nil // Clear before load to prevent dublicates
	err := c.Load()
	if err != nil {
		ctxLog.Fatalf("Config file '%v' is invalid or wrong file. '%v'", c.ConfigFile, err)
	}

	ctxLog.Tracef("Config struct: '%+v'", c)
	ctxLog.Tracef("Config streams dir: '%v'", c.StreamsDir)
	return c
}

func (c *Config) AutomaticEnv(ctxLog *log.Entry) {
	c.SetLogLevel(ctxLog)
	c.SetConfFile(ctxLog)
	c.SetStreamsDir(ctxLog)
}

func (c *Config) SetLogLevel(ctxLog *log.Entry) {
	if FlagLogLevel != "" {
		switch FlagLogLevel {
		case "trace":
			ctxLog.Info("Log level: trace")
			log.SetLevel(log.TraceLevel)
		case "debug":
			ctxLog.Info("Log level: debug")
			log.SetLevel(log.DebugLevel)
		default:
			ctxLog.Info("Log level: info (default)")
			log.SetLevel(log.InfoLevel)
		}
	}

	if FlagLogLevel == "" && EnvLogLevel != "" {
		switch EnvLogLevel {
		case "trace":
			ctxLog.Info("Log level: trace")
			log.SetLevel(log.TraceLevel)
		case "debug":
			ctxLog.Info("Log level: debug")
			log.SetLevel(log.DebugLevel)
		default:
			ctxLog.Info("Log level: info (default)")
			log.SetLevel(log.InfoLevel)
		}
	}
	if FlagLogLevel == "" && EnvLogLevel == "" {
		ctxLog.Info("Log level: info (default)")
		log.SetLevel(log.InfoLevel)
	}
}

func (c *Config) SetConfFile(ctxLog *log.Entry) {
	// Find OS-specific home directory. Prepare configDir, configFilePath.
	homeDir, err := homedir.Dir()
	cobra.CheckErr(err)
	ctxLog.Tracef("User home directory: %s", homeDir)

	configDir = homeDir + sep + "rekoda"
	c.ConfigDir = configDir

	streamsDir = c.ConfigDir + sep + "streams"
	c.StreamsDir = streamsDir

	// Scenario: custom conf file via --config flag
	if FlagConfigFile != "" {
		ctxLog.Infof("Using custom --config file: '%v'", FlagConfigFile)
		c.ConfigFile = FlagConfigFile
		if exists := c.ConfigFileExists(c.ConfigFile); !exists {
			ctxLog.Fatal("Config file you stated does not exists. Double check your --config flag file path to rekoda.toml")
		}
	}

	// Scenario: custom conf file via 'REKODA_CONF_FILE' env
	if FlagConfigFile == "" && EnvConfigFile != "" {
		c.ConfigFile = EnvConfigFile
		ctxLog.Infof("Using REKODA_CONF_FILE environment: '%v'", EnvConfigFile)
		if exists := c.ConfigFileExists(EnvConfigFile); !exists {
			ctxLog.Fatalf("Config file '%v' does not exists. Double check your REKODA_CONF_FILE environement, or unset it.", EnvConfigFile)
		}
	}

	// Scenario: Using default config file path when no ENVs or custom flag stated
	if FlagConfigFile == "" && EnvConfigFile == "" {

		ctxLog.Tracef("Config streams dir: '%v'", c.StreamsDir)

		c.ConfigFile = c.ConfigDir + sep + ConfigFile
		ctxLog.Infof("Using default config file path: %s", c.ConfigFile)
		if exists := c.ConfigFileExists(c.ConfigFile); !exists {
			ctxLog.Infof("Config file: %s not found. First time running? Trying to create default config file", c.ConfigFile)
			err := c.CreateDefaultConf()
			if err != nil {
				ctxLog.Errorf("failed to create default configuration: %v", err)
			}

			ctxLog.Infof("Config file '%s' created", c.ConfigFile)
		}
	}
}

func (c *Config) SetStreamsDir(ctxLog *log.Entry) {
	// Set custom output dir if stated
	if FlagStreamsDir != "" {
		ctxLog.Infof("Using custom --output folder: '%v'", FlagStreamsDir)
		c.StreamsDir = FlagStreamsDir
	}

	// Set custom conf dir via 'REKODA_STREAMS_DIR' env
	if FlagStreamsDir == "" && EnvStreamsDir != "" {
		c.StreamsDir = EnvStreamsDir
		ctxLog.Infof("Using REKODA_STREAMS_DIR environment: '%v'", EnvStreamsDir)
	}
}

// Save writes current Config struct into rekoda.toml file
func (c *Config) Save() error {
	// log.Println(ConfigStruct.Channels)
	if err := conf.Save(c.ConfigFile, &c); err != nil {
		return err
	}
	// conf.Load(configFilePath, &ConfigStruct)
	return nil
}

// Load reads back from rekoda.toml file and unmarshal int Config struct
func (c *Config) Load() error {
	// log.Printf("ConfigFilePath: '%v' \n c: '%v'", ConfigFilePath, c)
	err := conf.Load(c.ConfigFile, c)
	if err != nil {
		return err
	}
	return nil
}

// ConfigFileExists checks if config file exists in the stated filepath parameter
func (c *Config) ConfigFileExists(filepath string) bool {
	if _, err := os.Stat(filepath); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

// CreateDefaultConf creates default rekoda.toml conf file and writes it to disk.
// Used by default when --config or REKODA_CONF_FILE environement is not specified by user and rekoda.toml file does not exists in its default path ($HOME/rekoda/rekoda.toml)
func (c *Config) CreateDefaultConf() error {
	ctxLog := log.WithField("general", "MAKE CONF")

	// Make conf directory
	if err := os.MkdirAll(c.ConfigDir, 0777); err != nil {
		return err
	}
	ctxLog.Tracef("Config directory '%s' created", c.ConfigDir)

	// Make conf file
	c.MakeDefaultConfStruct()

	ctxLog.Tracef("Trying to write config to file '%v'", c.ConfigFile)
	if err := c.Save(); err != nil {
		return err
	}
	ctxLog.Tracef("%v 📁  wrote!", c.ConfigFile)

	return nil
}

// IsChannelInConfig checks if specific channel is in Config struct
func (c *Config) IsChannelInConfig(channel string) bool {
	for _, u := range c.Channels {
		if u.User == channel {
			return true
		}
	}
	return false
}

func ChannelAdd(username string) {
	// NYI check if twitch username exists/banned and then get ID

}

func ChannelUpdateID(username string, id int64) {

}

func ChannelEnable(username string) {

}

func ChannelDisable(username string) {

}

func ChannelDelete(username string) {

}

// ChannelList lists all channels in Config struct
func (c *Config) ChannelList() (string, string) {
	totalAmount := len(c.Channels)
	recAmount := 0
	var listOfAllChannels, listOfRecChannels string
	for _, v := range c.Channels {
		listOfAllChannels += fmt.Sprintf("%v ", v.User)
		if v.Enabled {
			recAmount++
			listOfRecChannels += fmt.Sprintf("%v ", v.User)
		}
	}
	total := fmt.Sprintf("%v channel(s) total in config: %v", totalAmount, listOfAllChannels)
	rec := fmt.Sprintf("%v channel(s) watching for recording: %v\n", recAmount, listOfRecChannels)
	return total, rec
}

// GetRecAmountChannels returns amount of channels being enabled for record
func (c *Config) GetRecAmountChannels() int {
	i := 0
	for _, c := range c.Channels {
		if c.Enabled {
			i++
		}
	}
	return i
}
