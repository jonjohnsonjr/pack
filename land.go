package pack

import (
	"context"

	"github.com/buildpack/lifecycle"
	"github.com/buildpack/lifecycle/metadata"
	"github.com/pkg/errors"

	"github.com/buildpack/pack/builder"
	"github.com/buildpack/pack/style"
)

type LandOptions struct {
	RepoName          string
	Publish           bool
	SkipPull          bool
	RunImage          string
	AdditionalMirrors map[string][]string
}

func (c *Client) Land(ctx context.Context, opts LandOptions) error {
	imageRef, err := c.parseTagReference(opts.RepoName)
	if err != nil {
		return errors.Wrapf(err, "invalid image name '%s'", opts.RepoName)
	}

	appImage, err := c.imageFetcher.Fetch(ctx, opts.RepoName, !opts.Publish, !opts.SkipPull)
	if err != nil {
		return err
	}

	md, err := metadata.GetLayersMetadata(appImage)
	if err != nil {
		return err
	}

	runImageName := c.resolveRunImage(
		opts.RunImage,
		imageRef.Context().RegistryStr(),
		builder.StackMetadata{
			RunImage: builder.RunImageMetadata{
				Image:   md.Stack.RunImage.Image,
				Mirrors: md.Stack.RunImage.Mirrors,
			},
		},
		opts.AdditionalMirrors)

	if runImageName == "" {
		return errors.New("run image must be specified")
	}

	baseImage, err := c.imageFetcher.Fetch(ctx, runImageName, !opts.Publish, !opts.SkipPull)
	if err != nil {
		return err
	}

	c.logger.Infof("Rebasing %s on run image %s", style.Symbol(appImage.Name()), style.Symbol(baseImage.Name()))
	lander := &lifecycle.Lander{Logger: c.logger}
	err = lander.Land(appImage, baseImage, nil)
	if err != nil {
		return err
	}

	appImageIdentifier, err := appImage.Identifier()
	if err != nil {
		return err
	}

	c.logger.Infof("Landed Image: %s", style.Symbol(appImageIdentifier.String()))
	return nil
}
