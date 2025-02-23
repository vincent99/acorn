package cli

import (
	"fmt"
	"os"

	cli "github.com/acorn-io/runtime/pkg/cli/builder"
	"github.com/acorn-io/runtime/pkg/client"
	acornsign "github.com/acorn-io/runtime/pkg/cosign"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

func NewImageVerify(c CommandContext) *cobra.Command {
	cmd := cli.Command(&ImageVerify{client: c.ClientFactory}, cobra.Command{
		Use: "verify IMAGE_NAME [flags]",
		Example: `# Verify using a locally stored public key file
acorn image verify my-image --key ./my-key.pub

# Verify using a public key belonging to a GitHub Identity
acorn image verify my-image --key gh://ibuildthecloud

# Verify using a public key belonging to an Acorn Manager Identity
acorn image verify my-image --key ac://ibuildthecloud
`,
		SilenceUsage:      true,
		Short:             "Verify Image Signatures",
		ValidArgsFunction: newCompletion(c.ClientFactory, imagesCompletion(true)).complete,
		Args:              cobra.ExactArgs(1),
	})
	_ = cmd.MarkFlagFilename("key")
	return cmd
}

type ImageVerify struct {
	client      ClientFactory
	Key         string            `usage:"Key to use for verifying" short:"k" local:"true"`
	Annotations map[string]string `usage:"Annotations to check for in the signature" short:"a" local:"true" name:"annotation"`
}

func (a *ImageVerify) Run(cmd *cobra.Command, args []string) error {
	if a.Key == "" {
		return fmt.Errorf("key is required")
	}

	imageName := args[0]

	c, err := a.client.CreateDefault()
	if err != nil {
		return err
	}

	auth, err := getAuthForImage(cmd.Context(), a.client, imageName)
	if err != nil {
		return err
	}

	// not failing here, since it could be a local image
	ref, _ := name.ParseReference(imageName)

	details, err := c.ImageDetails(cmd.Context(), args[0], &client.ImageDetailsOptions{
		Auth: auth,
	})
	if err != nil {
		return err
	}

	targetDigest := ref.Context().Digest(details.AppImage.Digest)

	pterm.Info.Printf("Verifying Image %s (digest: %s) using key %s\n", imageName, targetDigest, a.Key)

	vOpts := &client.ImageVerifyOptions{
		Annotations: a.Annotations,
		PublicKey:   a.Key,
		Auth:        auth,
	}

	// load public key from file (if it is a file, not a remote reference)
	if _, err := os.Stat(a.Key); err == nil {
		keyFileBytes, err := os.ReadFile(a.Key)
		if err != nil {
			return err
		}

		verifiers, err := acornsign.LoadVerifiers(cmd.Context(), string(keyFileBytes), "sha256")
		if err != nil {
			return err
		}

		pubkey, err := verifiers[0].PublicKey()
		if err != nil {
			return err
		}
		pem, _, err := acornsign.PemEncodeCryptoPublicKey(pubkey)
		if err != nil {
			return err
		}
		vOpts.PublicKey = string(pem)
	}

	_, err = c.ImageVerify(cmd.Context(), imageName, vOpts)
	if err == nil {
		pterm.Success.Println("Signature verified")
	}

	return err
}
