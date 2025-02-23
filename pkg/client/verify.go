package client

import (
	"context"
	"fmt"
	"strings"

	apiv1 "github.com/acorn-io/runtime/pkg/apis/api.acorn.io/v1"
	internalv1 "github.com/acorn-io/runtime/pkg/apis/internal.acorn.io/v1"
)

func (c *DefaultClient) ImageVerify(ctx context.Context, image string, opts *ImageVerifyOptions) (*apiv1.ImageSignature, error) {
	sigInput := &apiv1.ImageSignature{
		PublicKey: opts.PublicKey,
		Auth:      opts.Auth,
	}

	if opts.PublicKey == "" {
		return nil, fmt.Errorf("public key required for verification")
	}

	sigInput.Annotations = internalv1.SignatureAnnotations{
		Match: opts.Annotations,
	}

	imageDetails, err := c.ImageDetails(ctx, image, &ImageDetailsOptions{Auth: opts.Auth})
	if err != nil {
		return nil, err
	}

	image = strings.ReplaceAll(imageDetails.AppImage.ID, "/", "+")

	sigResult := &apiv1.ImageSignature{}
	err = c.RESTClient.Post().
		Namespace(c.Namespace).
		Resource("images").
		Name(image).
		SubResource("verify").
		Body(sigInput).Do(ctx).Into(sigResult)

	return sigResult, err
}
