package cmd

import (
	"testing"
)

func TestVersionCommand(t *testing.T) {
	cmd := NewVersionCmd()
	_, err := ExecuteCommand(cmd)
	if err != nil {
		t.Error(err)
	}
}
