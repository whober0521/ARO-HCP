package main

import (
	"log"

	"github.com/Azure/ARO-HCP/tooling/poc"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func DefaultGenerationOptions() *GenerationOptions {
	return &GenerationOptions{}
}

type GenerationOptions struct {
	Region string
	User   string
}

func BindGenerationOptions(opts *GenerationOptions, flags *pflag.FlagSet) {
	flags.StringVar(&opts.Region, "region", opts.Region, "resources location")
	flags.StringVar(&opts.User, "user", opts.User, "unique user name")
}

func main() {
	cmd := &cobra.Command{}

	opts := DefaultGenerationOptions()
	BindGenerationOptions(opts, cmd.Flags())
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		println("Region:", opts.Region)
		println("User:", opts.User)
		println()

		if err := poc.Print(opts.Region, opts.User); err != nil {
			return err
		}

		return nil
	}

	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
