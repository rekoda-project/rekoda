package main

import (
	"os"

	nested "github.com/antonfisher/nested-logrus-formatter"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&nested.Formatter{
		TimestampFormat: "[2006 Jan 2, Monday][15:04:05 MST]",
		HideKeys:        true,
		FieldsOrder:     []string{"general", "func", "status", "channel", "file"},
	})
	// Output to stdout instead of the default stderr
	log.SetOutput(os.Stdout)

	log.SetLevel(log.InfoLevel)
	//	log.Infof("Rekoda started")
}
