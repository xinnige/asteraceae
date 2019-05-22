package awsapi

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/xinnige/asteraceae/calendula/mock"
)

func TestNewSSMAPI(t *testing.T) {
	sess := &AWSServiceSession{}
	ssmapi := NewSSMAPI(sess)
	assert.NotNil(t, ssmapi)
}

func TestNewSSMAPIError(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockSession := mock.NewMockAWSSession(mockCtrl)
	mockSession.EXPECT().NewSession().Return(
		nil, errors.New("FakeNewSessionError")).Times(1)

	ssmapi := NewSSMAPI(mockSession)
	assert.Nil(t, ssmapi)
}

func TestPutParameter(t *testing.T) {
	awsapi := AWSAPI{}
	ssmapi := mock.SSMPutParameter(32)
	tags := map[string]string{"key1": "val1", "key2": "val2"}

	assert.Nil(t, awsapi.PutParameter(ssmapi, "fake-name", "a,b,c", "", tags, false, true))
	assert.Nil(t, awsapi.PutParameter(ssmapi, "fake-name", "fakevalue", "", tags, false, false))
	assert.Nil(t, awsapi.PutParameter(ssmapi, "fake-name", "fakevalue", "fake-key", nil, true, false))
}

func TestPutParameterError(t *testing.T) {
	awsapi := AWSAPI{}
	ssmapi := mock.SSMPutParameterError()
	err := awsapi.PutParameter(ssmapi, "fake-name", "fakevalue", "", nil, true, false)
	assert.NotNil(t, err)
	assert.Equal(t, "FakeSSMPutError", err.Error())

	err = awsapi.PutParameter(ssmapi, "fake-name", "fakevalue", "fake-key", nil, true, false)
	assert.NotNil(t, err)
	assert.Equal(t, "FakeSSMPutError", err.Error())
}

func TestAddTags(t *testing.T) {
	awsapi := AWSAPI{}
	ssmapi := mock.SSMAddTags()
	tags := map[string]string{"key1": "val1", "key2": "val2"}
	assert.Nil(t, awsapi.AddTagsToResource(ssmapi, "fake-name", "fake-type", tags))
}

func TestAddTagsError(t *testing.T) {
	awsapi := AWSAPI{}
	ssmapi := mock.SSMAddTagsError()
	tags := map[string]string{"key1": "val1", "key2": "val2"}
	err := awsapi.AddTagsToResource(ssmapi, "fake-name", "fake-type", tags)
	assert.Equal(t, "FakeSSMAddTagsError", err.Error())
}

func TestGetParams(t *testing.T) {
	awsapi := AWSAPI{}
	names := [][]string{
		[]string{"n11", "n12"}, []string{"n21", "n22", "n23"}, []string{"n31"}}
	values := [][]string{
		[]string{"v11", "v12"}, []string{"v21", "v22", "v23"}, []string{"v31"}}
	nexts := []*string{aws.String("fake-next-1"), aws.String("fake-next-2"), nil}
	ssmapi := mock.SSMGetParams(names, values, nexts)
	result, err := awsapi.GetParametersByPath(ssmapi, "fake-name", true, 10)
	assert.Equal(t, 6, len(result))
	assert.Nil(t, err)

	empty := [][]string{}
	ssmapi = mock.SSMGetParams(empty, empty, []*string{})
	result, err = awsapi.GetParametersByPath(ssmapi, "fake-name", true, 10)
	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.Equal(t, "parameter fake-name not found", err.Error())
}

func TestGetParamsError(t *testing.T) {
	awsapi := AWSAPI{}
	ssmapi := mock.SSMGetParamsError()
	_, _, err := awsapi.GetParametersByPathIter(ssmapi,
		"fake-name", aws.String("fake-next"), true, 10)
	assert.Equal(t, "FakeSSMAGetParamsError", err.Error())

	_, err = awsapi.GetParametersByPath(ssmapi, "fake-name", true, 10)
	assert.Equal(t, "FakeSSMAGetParamsError", err.Error())
}

func TestListTags(t *testing.T) {
	awsapi := AWSAPI{}
	names := [][]string{
		[]string{"n11", "n12"}, []string{"n21", "n22", "n23"}, []string{"n31"}}
	values := [][]string{
		[]string{"v11", "v12"}, []string{"v21", "v22", "v23"}, []string{"v31"}}
	ssmapi := mock.SSMListTags(names, values)
	result, err := awsapi.ListTagsForResource(ssmapi, "fake-name", "fake-type")
	assert.Equal(t, 2, len(result))
	assert.Nil(t, err)
	result, err = awsapi.ListTagsForResource(ssmapi, "fake-name", "fake-type")
	assert.Equal(t, 3, len(result))
	assert.Nil(t, err)
	result, err = awsapi.ListTagsForResource(ssmapi, "fake-name", "fake-type")
	assert.Equal(t, 1, len(result))
	assert.Nil(t, err)
}

func TestListTagsError(t *testing.T) {
	awsapi := AWSAPI{}
	ssmapi := mock.SSMListTagsError()
	_, err := awsapi.ListTagsForResource(ssmapi, "fake-name", "fake-type")
	assert.Equal(t, "FakeSSMListTagsError", err.Error())
}

func TestDeleteParams(t *testing.T) {
	awsapi := AWSAPI{}
	ssmapi := mock.SSMDeleteParams([]string{"abc"}, []string{"123"})
	err := awsapi.DeleteParameters(ssmapi, []string{"abc", "123"})
	assert.Nil(t, err)
}

func TestDeleteParamsError(t *testing.T) {
	awsapi := AWSAPI{}
	ssmapi := mock.SSMDeleteParamsError()
	err := awsapi.DeleteParameters(ssmapi, []string{})
	assert.Equal(t, "FakeSSMDeleteParamsError", err.Error())
}
