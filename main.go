package main

import (
	log "github.com/sirupsen/logrus"

	"github.com/ankittk/go-dependency-tree-parser/cmd"
)

func main() {
	cmd.Execute()
}

func init() {
	log.SetFormatter(&log.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})
	log.SetLevel(log.InfoLevel)
}
