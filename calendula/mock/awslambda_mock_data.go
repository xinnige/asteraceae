package mock

import (
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/aws/aws-sdk-go/service/lambda/lambdaiface"
)

// LambdaAPI mocks s3iface.S3API
type LambdaAPI struct {
	lambdaiface.LambdaAPI
	invokeOutput  lambda.InvokeOutput
	addpermOutput lambda.AddPermissionOutput
	invokeErr     error
	addpermErr    error
}

// Invoke mocks lambda.Invoke
func (m LambdaAPI) Invoke(
	input *lambda.InvokeInput) (*lambda.InvokeOutput, error) {
	return &m.invokeOutput, m.invokeErr
}

// AddPermission mocks lambda.AddPermission
func (m LambdaAPI) AddPermission(
	input *lambda.AddPermissionInput) (*lambda.AddPermissionOutput, error) {
	return &m.addpermOutput, m.addpermErr
}

func mockInvokeOutput(
	payload []byte, status int64, errString string) *lambda.InvokeOutput {
	return &lambda.InvokeOutput{
		FunctionError: aws.String(errString),
		Payload:       payload,
		StatusCode:    aws.Int64(status),
	}
}

func mockAddPermOutput(state string) *lambda.AddPermissionOutput {
	return &lambda.AddPermissionOutput{
		Statement: aws.String(state),
	}
}

// LambdaAPIInvoke mocks a Lambda API
func LambdaAPIInvoke() *LambdaAPI {
	return &LambdaAPI{
		invokeOutput: *mockInvokeOutput([]byte("payload"), 202, ""),
		invokeErr:    nil,
	}
}

// LambdaAPIInvokeError mocks a Lambda API w/ error
func LambdaAPIInvokeError() *LambdaAPI {
	return &LambdaAPI{
		invokeOutput: *mockInvokeOutput([]byte("payload"), 500, ""),
		invokeErr:    errors.New("MockInvokeError"),
	}
}

// LambdaAPIAddPerm mocks a Lambda.AddPermission API
func LambdaAPIAddPerm(state string) *LambdaAPI {
	return &LambdaAPI{
		addpermOutput: *mockAddPermOutput(state),
		addpermErr:    nil,
	}
}

// LambdaAPIAddPermError mocks a Lambda.AddPermission API w/ error
func LambdaAPIAddPermError() *LambdaAPI {
	return &LambdaAPI{
		addpermOutput: *mockAddPermOutput(""),
		addpermErr:    errors.New("MockAddPermError"),
	}
}
