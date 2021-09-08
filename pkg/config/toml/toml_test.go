package toml

import (
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var filepath string = "test.toml"

type Config struct {
	Title string `toml:"title"`
}

func TestSave(t *testing.T) {
	c := &Config{Title: "test"}
	err := Save(filepath, &c)
	if err != nil {
		assert.Error(t, err)
	}

	tomlFile, err := os.Open(filepath)
	if err != nil {
		assert.Error(t, err)
	}
	defer tomlFile.Close()

	bytes, _ := io.ReadAll(tomlFile)
	assert.Equal(t, "title = 'test'\n", string(bytes))
}

func TestLoad(t *testing.T) {
	c := &Config{}
	err := Load(filepath, &c)
	if err != nil {
		assert.Error(t, err)
	}
	assert.Equal(t, "test", c.Title)
}
