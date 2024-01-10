package main

import (
	"github.com/marcuzh/ansible-superputty/cmd"
	log "github.com/sirupsen/logrus"
)

func main() {
	rootCmd, err := cmd.NewRootCmd()
	if err != nil {
		log.Fatal(err)
	}

	rootCmd.Setup()

	err = rootCmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
