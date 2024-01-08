package main

import (
	"github.com/marcuzh/ansible-superputty/cmd"
	log "github.com/sirupsen/logrus"
)

func main() {
	rootCmd := cmd.NewCmd()
	err := rootCmd.Setup()
	if err != nil {
		log.Fatal(err)
	}

	err = rootCmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
