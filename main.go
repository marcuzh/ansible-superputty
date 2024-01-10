package main

import (
	"github.com/marcuzh/ansible-superputty/cmd"
	log "github.com/sirupsen/logrus"
)

func main() {
	rootCmd := cmd.NewRootCmd()

	err := rootCmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
