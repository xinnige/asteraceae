package awsapi

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/xinnige/asteraceae/calendula/mock"
)

func TestNewLambdaAPI(t *testing.T) {
	sess := &AWSServiceSession{}
	lambdaapi := NewLambdaAPI(sess)
	assert.NotNil(t, lambdaapi)
}

func TestNewLambdaAPIError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockSession := mock.NewMockAWSSession(mockCtrl)
	mockSession.EXPECT().NewSession().Return(
		nil, errors.New("FakeNewSessionError")).Times(1)

	lambdaapi := NewLambdaAPI(mockSession)
	assert.Nil(t, lambdaapi)
}

func TestInvokeLambda(t *testing.T) {
	lambdaapi := mock.LambdaAPIInvoke()

	awsapi := AWSAPI{}
	output, err := awsapi.InvokeLambda(lambdaapi, "fakearn", []byte("json"), "Event")
	assert.NotNil(t, output)
	assert.Nil(t, err)

	lambdaapi = mock.LambdaAPIInvokeError()
	output, err = awsapi.InvokeLambda(lambdaapi, "fakearn", []byte("json"), "Event")
	assert.Nil(t, output)
	assert.NotNil(t, err)
	assert.Equal(t, "MockInvokeError", err.Error())

}

func TestLambdaAddPermission(t *testing.T) {

	lambdasvc := mock.LambdaAPIAddPerm("fake-state")

	awsapi := AWSAPI{}
	err := awsapi.LambdaAddPermission(lambdasvc, ActionInvoke,
		"lambda_function_name", PrincipalEvents,
		"arn:aws:events:ap-northeast-1:xxxxx:rule/lambdaname-randid",
		"randstateid12345")

	assert.Nil(t, err)

	lambdasvc = mock.LambdaAPIAddPermError()
	err = awsapi.LambdaAddPermission(lambdasvc, ActionInvoke,
		"lambda_function_name", PrincipalEvents,
		"arn:aws:events:ap-northeast-1:xxxxx:rule/lambdaname-randid",
		"randstateid12345")
	assert.Equal(t, "MockAddPermError", err.Error())
}
