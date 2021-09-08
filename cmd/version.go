package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const (
	APP     = "Rekōdā"
	VERSION = "0.0.1"
)

func NewVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "version",
		Version: VERSION,
		Short:   "Prints version",
		RunE: func(cmd *cobra.Command, args []string) error {
			log.WithField("general", "CLI").Infof("%v v%v\n", APP, VERSION)
			return nil
		},
	}
}

var versionCmd = NewVersionCmd()

func init() {
	rootCmd.AddCommand(versionCmd)
}
