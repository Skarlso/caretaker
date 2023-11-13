package main

import (
	"log"

	"github.com/skarlso/caretaker/cmd"
)

func main() {
	root := cmd.CreateRootCommand()
	if err := root.Execute(); err != nil {
		log.Fatal(err)
	}
}
