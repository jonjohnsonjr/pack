package commands

import (
	"github.com/spf13/cobra"

	"github.com/buildpack/pack"
	"github.com/buildpack/pack/config"
	"github.com/buildpack/pack/logging"
	"github.com/buildpack/pack/style"
)

func Land(logger logging.Logger, cfg config.Config, client PackClient) *cobra.Command {
	var opts pack.LandOptions
	ctx := createCancellableContext()

	cmd := &cobra.Command{
		Use:   "land <image-name>",
		Args:  cobra.ExactArgs(1),
		Short: "Land an image by removing the launcher",
		RunE: logError(logger, func(cmd *cobra.Command, args []string) error {
			opts.RepoName = args[0]
			opts.AdditionalMirrors = getMirrors(cfg)
			if err := client.Land(ctx, opts); err != nil {
				return err
			}
			logger.Infof("Successfully landed image %s", style.Symbol(opts.RepoName))
			return nil
		}),
	}
	cmd.Flags().BoolVar(&opts.Publish, "publish", false, "Publish to registry")
	cmd.Flags().BoolVar(&opts.SkipPull, "no-pull", false, "Skip pulling app and run images before use")
	AddHelpFlag(cmd, "land")
	return cmd
}
