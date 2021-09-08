package config

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	//"github.com/wmw9/rekoda/internal/config"
)

func TestMakeDefaultConfStruct(t *testing.T) {
	c := InitConfig()
	c.MakeDefaultConfStruct()
	want := &Config{
		Title:      "Rekoda configuration file",
		Version:    1,
		StreamsDir: "streams",
		Channels: []Channels{
			{
				Enabled: false,
				User:    "reckful",
				ID:      999999999999,
				Quality: "best",
			},
		},
	}
	if !cmp.Equal(c, want) {
		t.Errorf("got: '%v'\n wanted: '%v'", c, want)
	}

}

func TestInitConfig(t *testing.T) {
	want := "Rekoda configuration file"
	c := InitConfig()
	c.MakeDefaultConfStruct()
	assert.Equal(t, want, c.Title)
}

func TestSave(t *testing.T) {
	c := InitConfig()
	err := c.Save()
	if err != nil {
		assert.Error(t, err)
	}
}

func TestLoad(t *testing.T) {
	want := "Rekoda configuration file"
	c := InitConfig()
	c.Load()
	assert.Equal(t, c.Title, want)
}

func TestConfigFileExists(t *testing.T) {
	c := InitConfig()
	assert.Equal(t, c.ConfigFileExists(c.ConfigFile), true)
}

func TestIsChannelInConfig(t *testing.T) {
	c := InitConfig()
	assert.Equal(t, !c.IsChannelInConfig("reckful"), true)
}
