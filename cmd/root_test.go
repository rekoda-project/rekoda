package cmd

import (
	"bytes"
	"testing"

	"github.com/spf13/cobra"
)

func ExecuteCommand(root *cobra.Command, args ...string) (output string, err error) {
	_, output, err = ExecuteCommandC(root, args...)
	return output, err
}

func ExecuteCommandC(root *cobra.Command, args ...string) (c *cobra.Command, output string, err error) {
	buf := new(bytes.Buffer)
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs(args)

	c, err = root.ExecuteC()

	return c, buf.String(), err
}

func TestCallCommandWithoutSubcommands(t *testing.T) {
	cmd := NewRootCmd()
	_, err := ExecuteCommand(cmd)
	if err != nil {
		t.Errorf("Calling command without subcommands should not have error: %v", err)
	}
}
