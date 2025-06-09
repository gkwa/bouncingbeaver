package main

import (
	"log"

	"github.com/gkwa/bouncingbeaver/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
