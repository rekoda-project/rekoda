package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewChannelCmd(t *testing.T) {
	rootCmd := NewRootCmd()
	channelCmd := NewChannelCmd()

	rootCmd.AddCommand(channelCmd)
	c, out, err := ExecuteCommandC(rootCmd, "channel")
	if err != nil {
		t.Error(err, out)
	}
	assert.Equal(t, "channel", c.Name())

}

func TestNewAddCmd(t *testing.T) {
	channelCmd := NewChannelCmd()
	addCmd := NewAddCmd()

	channelCmd.AddCommand(addCmd)

	c, _, err := ExecuteCommandC(channelCmd, "add", "reckful")
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "add", c.Name())
}

func TestNewRemoveCmd(t *testing.T) {
	channelCmd := NewChannelCmd()
	removeCmd := NewRemoveCmd()

	channelCmd.AddCommand(removeCmd)
	//rootCmd.AddCommand(channelCmd)

	c, _, err := ExecuteCommandC(channelCmd, "remove", "reckful")
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "remove", c.Name())
}

func TestNewListCmd(t *testing.T) {
	channelCmd := NewChannelCmd()
	listCmd := NewListCmd()

	channelCmd.AddCommand(listCmd)
	c, out, err := ExecuteCommandC(channelCmd, "list")
	if err != nil {
		t.Error(err, out)
	}
	assert.Equal(t, "list", c.Name())
}

func TestDisableCmd(t *testing.T) {
	channelCmd := NewChannelCmd()
	disableCmd := NewDisableCmd()

	channelCmd.AddCommand(disableCmd)
	c, out, err := ExecuteCommandC(channelCmd, "disable", "reckful")
	if err != nil {
		t.Error(err, out)
	}
	assert.Equal(t, "disable", c.Name())
}

func TestEnableCmd(t *testing.T) {
	channelCmd := NewChannelCmd()
	enableCmd := NewEnableCmd()

	channelCmd.AddCommand(enableCmd)
	c, out, err := ExecuteCommandC(channelCmd, "enable", "reckful")
	if err != nil {
		t.Error(err, out)
	}
	assert.Equal(t, "enable", c.Name())
}
