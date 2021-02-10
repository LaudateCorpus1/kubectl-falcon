package cmd

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/containers/image/v5/copy"
	"github.com/containers/image/v5/signature"
	"github.com/containers/image/v5/transports/alltransports"
	"github.com/containers/image/v5/types"
	"github.com/spf13/cobra"

	"github.com/crowdstrike/gofalcon/falcon"
	"github.com/crowdstrike/kubectl-falcon/pkg/docker_creds"
	"github.com/crowdstrike/kubectl-falcon/pkg/falcon_image"
)

type refreshImageOptions struct {
	global              *globalOptions
	destinationUsername string
}

func imageRefresh(global *globalOptions) *cobra.Command {
	opts := refreshImageOptions{
		global: global,
	}

	cmd := &cobra.Command{
		Use:   "image-refresh DESTINATION-IMAGE",
		Short: "Refresh falcon container image in the registry or local docker daemon",
		Long:  "Fetches falcon container image from the falcon.crowdstrike.com and pushes it to the container registry or local docker daemon",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return errors.New("Please provide exactly one argument as DESTINATION-IMAGE")
			}
			if os.Getenv("FALCON_CLIENT_ID") == "" {
				return errors.New("Missing FALCON_CLIENT_ID environment variable. Please provide your Client Id for authentication with CrowdStrike Falcon platform")
			}
			if os.Getenv("FALCON_CLIENT_SECRET") == "" {
				return errors.New("Missing FALCON_CLIENT_SECRET environment variable. Please provide your Client Secret for authentication with CrowdStrike Falcon platform")
			}
			return nil
		},
		RunE: commandAction(opts.run),
		Example: `
# Push image directly from CrowdStrike Falcon API to remote registry
REGISTRY_PASSWORD=$(aws ecr get-login-password --region eu-west-1) kubectl falcon image-refresh -u AWS docker://123456.dkr.ecr.REGION.amazonaws.com/falcon-sensor:latest

# Push image from CrowdStrike Falcon API to local docker daemon
kubectl falcon image-refresh docker-daemon:falcon-sensor:latest
`,
	}
	cmd.Flags().StringVarP(&opts.destinationUsername, "dest-user", "u", "", "Username for accessing the push registry, use $REGISTRY_PASSWORD env variable to supply password")

	return cmd
}

func (opts *refreshImageOptions) run(args []string, stdout io.Writer) error {
	ctx, cancel := opts.global.commandTimeoutContext()
	defer cancel()

	policy := &signature.Policy{Default: []signature.PolicyRequirement{signature.NewPRInsecureAcceptAnything()}}
	policyContext, err := signature.NewPolicyContext(policy)
	if err != nil {
		return fmt.Errorf("Error loading trust policy: %v", err)
	}
	defer func() { _ = policyContext.Destroy() }()

	destRef, err := alltransports.ParseImageName(args[0])
	if err != nil {
		return fmt.Errorf("Invalid destination name %s: %v", args[0], err)
	}

	destinationContext, err := opts.destinationContext(destRef)
	if err != nil {
		return err
	}

	falconImage, err := falcon_image.Pull(&falcon.ApiConfig{
		ClientId:     os.Getenv("FALCON_CLIENT_ID"),
		ClientSecret: os.Getenv("FALCON_CLIENT_SECRET"),
		Context:      ctx}, stdout)
	if err != nil {
		return fmt.Errorf("Failed to pull falcon image: %v", err)
	}
	defer func() { 	_ = falconImage.Delete() }()

	ref, err := falconImage.ImageReference()
	if err != nil {
		return fmt.Errorf("Failed to build internal image representation for falcon image: %v", err)
	}

	fmt.Fprintf(stdout, "Pushing image to %s\n", destRef)
	_, err = copy.Image(ctx, policyContext, destRef, ref, &copy.Options{
		DestinationCtx: destinationContext,
		ReportWriter:   stdout,
	})
	return wrapWithHint(err)
}

func (opts *refreshImageOptions) destinationContext(imageRef types.ImageReference) (*types.SystemContext, error) {
	ctx := &types.SystemContext{
		DockerAuthConfig: &types.DockerAuthConfig{},
	}
	if opts.destinationUsername != "" {
		ctx.DockerAuthConfig.Username = opts.destinationUsername
		ctx.DockerAuthConfig.Password = os.Getenv("REGISTRY_PASSWORD")
	} else {
		auth, err := docker_creds.BestGuess(imageRef)
		if err != nil {
			return nil, err
		}
		ctx.DockerAuthConfig = auth
	}
	return ctx, nil

}

func wrapWithHint(in error) error {
	// Use of credentials store outside of docker command is somewhat limited
	// See https://github.com/moby/moby/issues/39377
	// https://github.com/containers/image/pull/656
	if in == nil {
		return in
	}
	if strings.Contains(in.Error(), "authentication required") {
		return fmt.Errorf(`%s
Could not find local authentication config for the registry.
Please use docker login, podman login, skopeo login, or other standard method of authenticating`, in)
	}
	return in
}
