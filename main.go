//go:generate wire ./service
//go:generate wire ./app

package main

import (
	"appsrv/cmd"
	"log"
)

func main() {
	if err := cmd.RootCommand.Execute(); err != nil {
		log.Fatal(err)
	}
}
