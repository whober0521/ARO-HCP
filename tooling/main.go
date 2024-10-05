package main

import (
	"log"

	"github.com/Azure/ARO-HCP/tooling/poc"
	"github.com/spf13/cobra"
)

func main() {
	var cmd = &cobra.Command{}

	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}

	if err := poc.Print("poc/config.yaml"); err != nil {
		log.Fatal(err)
	}
}
