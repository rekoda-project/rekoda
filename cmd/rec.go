package cmd

import (
	"github.com/spf13/cobra"
	"github.com/kappaflow/rekoda/internal/recorder"
)

// recCmd represents the rec command
func NewRecCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "rec",
		Short: "Start recording streams",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			recorder.Start()
		},
	}
}

var recCmd = NewRecCmd()

func init() {
	rootCmd.AddCommand(recCmd)
}
