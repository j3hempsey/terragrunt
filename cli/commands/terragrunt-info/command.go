package terragruntinfo

import (
	"github.com/gruntwork-io/terragrunt/options"
	"github.com/gruntwork-io/terragrunt/pkg/cli"
)

const (
	CommandName = "terragrunt-info"
)

func NewCommand(opts *options.TerragruntOptions) *cli.Command {
	command := &cli.Command{
		Name:  CommandName,
		Usage: "Emits limited terragrunt state on stdout and exits.",
		Action: func(ctx *cli.Context) error {
			if err := opts.InitialSetup(ctx); err != nil {
				return err
			}

			return Run(opts)
		},
	}

	return command
}
