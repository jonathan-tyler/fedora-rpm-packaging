package buildcontainer

import (
	"context"
	"fmt"

	"github.com/your-github-username/fedora-package-infra/builder/src/core/packages"
)

type ResolvedTag struct {
	PackageName string
	UpstreamTag string
	ArtifactTag string
}

func (s Service) ResolveArtifactTag(ctx context.Context, packageName string) (ResolvedTag, error) {
	definition, err := s.Registry.Lookup(packageName)
	if err != nil {
		return ResolvedTag{}, err
	}

	upstreamTag, err := s.resolveLatestStableTag(ctx, definition)
	if err != nil {
		return ResolvedTag{}, err
	}
	if upstreamTag == "" {
		return ResolvedTag{}, fmt.Errorf("no stable version tag found for %s", packageName)
	}

	return ResolvedTag{
		PackageName: packageName,
		UpstreamTag: upstreamTag,
		ArtifactTag: artifactDirectoryName(upstreamTag, upstreamTag),
	}, nil
}

func (s Service) resolveLatestStableTag(ctx context.Context, definition packages.Definition) (string, error) {
	output, err := s.Runner.Capture(ctx, s.Tools.Git.ListRemoteTagsCommand(definition.RepoURL, s.Stdout))
	if err != nil {
		return "", fmt.Errorf("resolve latest stable tag for %s: %w", definition.Name, err)
	}

	return latestStableVersionTag(output), nil
}
