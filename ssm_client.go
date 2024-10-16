package ssmdotenv

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/ssm"
)

// SSMClient defines the interface for SSM client methods used in the package.
type SSMClient interface {
	GetParameter(ctx context.Context, params *ssm.GetParameterInput, optFns ...func(*ssm.Options)) (*ssm.GetParameterOutput, error)
	GetParametersByPathPaginator(params *ssm.GetParametersByPathInput) *ssm.GetParametersByPathPaginator
}

type AWSSSMClient struct {
	client *ssm.Client
}

func (r *AWSSSMClient) GetParameter(ctx context.Context, params *ssm.GetParameterInput, optFns ...func(*ssm.Options)) (*ssm.GetParameterOutput, error) {
	return r.client.GetParameter(ctx, params, optFns...)
}

func (r *AWSSSMClient) GetParametersByPathPaginator(params *ssm.GetParametersByPathInput) *ssm.GetParametersByPathPaginator {
	return ssm.NewGetParametersByPathPaginator(r.client, params)
}
