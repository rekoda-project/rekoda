package toml

import (
	"bytes"
	"io"
	"os"

	"github.com/pelletier/go-toml/v2"
)

func Load(filepath string, config interface{}) error {
	tomlFile, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer tomlFile.Close()

	bytes, _ := io.ReadAll(tomlFile)

	// Unmarshal into global config variable
	err = toml.Unmarshal(bytes, config)
	if err != nil {
		return err
	}
	return nil
}

func Save(filepath string, config interface{}) error {
	buf := bytes.Buffer{}
	enc := toml.NewEncoder(&buf)
	enc.SetIndentTables(true)
	err := enc.Encode(config)
	if err != nil {
		return err
	}

	err = os.WriteFile(filepath, buf.Bytes(), 0777)
	if err != nil {
		return err
	}
	return nil
}

/*
func CreateDefault(filepath string, config interface{}) error {
	buf := bytes.Buffer{}
	enc := toml.NewEncoder(&buf)
	enc.SetIndentTables(true)
	err := enc.Encode(config)
	if err != nil {
		return err
	}

	err = os.WriteFile("123.toml", buf.Bytes(), 0777)
	if err != nil {
		panic(err)
	}
	log.Printf("%v 📁  saved.", tomlFile)
	return buf.Bytes(), nil
}
*/
