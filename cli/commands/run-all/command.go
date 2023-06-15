package runall

import (
	awsproviderpatch "github.com/gruntwork-io/terragrunt/cli/commands/aws-provider-patch"
	graphdependencies "github.com/gruntwork-io/terragrunt/cli/commands/graph-dependencies"
	"github.com/gruntwork-io/terragrunt/cli/commands/hclfmt"
	renderjson "github.com/gruntwork-io/terragrunt/cli/commands/render-json"
	"github.com/gruntwork-io/terragrunt/cli/commands/terraform"
	terragruntinfo "github.com/gruntwork-io/terragrunt/cli/commands/terragrunt-info"
	validateinputs "github.com/gruntwork-io/terragrunt/cli/commands/validate-inputs"
	"github.com/gruntwork-io/terragrunt/cli/flags"
	"github.com/gruntwork-io/terragrunt/options"
	"github.com/gruntwork-io/terragrunt/pkg/cli"
)

const (
	CommandName = "run-all"
)

var (
	TerragruntFlagNames = flags.CommonFlagNames
)

func NewCommand(opts *options.TerragruntOptions) *cli.Command {
	return &cli.Command{
		Name:        CommandName,
		Usage:       "Run a terraform command against a 'stack' by running the specified command in each subfolder.",
		Description: "The command will recursively find terragrunt modules in the current directory tree and run the terraform command in dependency order (unless the command is destroy, in which case the command is run in reverse dependency order).",
		Subcommands: cli.Commands{
			terragruntinfo.NewCommand(opts),    // terragrunt-info
			validateinputs.NewCommand(opts),    // validate-inputs
			graphdependencies.NewCommand(opts), // graph-dependencies
			hclfmt.NewCommand(opts),            // hclfmt
			renderjson.NewCommand(opts),        // render-json
			awsproviderpatch.NewCommand(opts),  // aws-provider-patch
			terraform.NewCommand(opts),         // *
		},
		Flags:  flags.NewFlags(opts).Filter(TerragruntFlagNames),
		Before: func(ctx *cli.Context) error { return ctx.App.Before(ctx) },
		Action: Action(opts),
	}
}

func Action(opts *options.TerragruntOptions) func(ctx *cli.Context) error {
	return func(ctx *cli.Context) error {
		opts.RunTerragrunt = func(opts *options.TerragruntOptions) error {
			for _, command := range ctx.Command.Subcommands {
				if command.Name == opts.TerraformCommand {
					ctx := ctx.WithValue(options.ContextKey, opts)
					return command.Action(ctx)
				}
			}

			return terraform.Run(opts)
		}

		return Run(opts.FromContext(ctx))
	}
}
