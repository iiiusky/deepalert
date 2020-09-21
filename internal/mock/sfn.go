package mock

import (
	"github.com/aws/aws-sdk-go/service/sfn"
	"github.com/deepalert/deepalert/internal/adaptor"
)

// NewSFnClient creates mock SNS client
func NewSFnClient(region string) (adaptor.SFnClient, error) {
	return &SFnClient{region: region}, nil
}

// SFnClient is mock
type SFnClient struct {
	region string
	input  []*sfn.StartExecutionInput
}

// StartExecution of mock SFnClient only stores sfn.StartExecutionInput
func (x *SFnClient) StartExecution(input *sfn.StartExecutionInput) (*sfn.StartExecutionOutput, error) {
	x.input = append(x.input, input)
	return &sfn.StartExecutionOutput{}, nil
}
