package awsapi

import (
	"errors"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/xinnige/asteraceae/calendula/mock"
)

func TestNewS3API(t *testing.T) {
	sess := &AWSServiceSession{}
	s3api := NewS3API(sess)
	assert.NotNil(t, s3api)
}

func TestNewS3APIError(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockSession := mock.NewMockAWSSession(mockCtrl)
	mockSession.EXPECT().NewSession().Return(
		nil, errors.New("FakeNewSessionError")).Times(1)

	s3api := NewS3API(mockSession)
	assert.Nil(t, s3api)
}

func TestNewUploaderAPI(t *testing.T) {
	sess := &AWSServiceSession{}
	s3api := NewS3API(sess)
	uploader := NewS3UploaderAPI(s3api)
	assert.NotNil(t, uploader)
}

func TestGetObject(t *testing.T) {
	bucket := "example-bucket"
	key := "example-key/example.json"
	awsapi := AWSAPI{}

	content := "{\"samplejson\":{}}"
	versionID := "ESGHby71m6xT4QoerT81suYiCzOod8v_"
	mockS3API := mock.S3APIGetObjects([]string{content}, versionID)
	result, err := awsapi.GetObject(mockS3API, bucket, key)
	assert.Nil(t, err)
	assert.Equal(t, content, result)
}

func TestGetObjectError(t *testing.T) {

	bucket := "example-bucket"
	key := "example-key/example.json"
	awsapi := AWSAPI{}

	mockS3API := mock.S3APIGetObjectError()
	result, err := awsapi.GetObject(mockS3API, bucket, key)
	assert.NotNil(t, err)
	assert.Equal(t, "", result)
	assert.Equal(t, "FakeGetObjectError", err.Error())
}

func TestListObjects(t *testing.T) {
	bucket := "example-bucket"
	prefix := "example-prefix/prefix/example"
	awsapi := AWSAPI{}

	// 1 time access for all
	names := []string{"name1", "name2"}
	mockS3API := mock.S3APIListObjects(bucket, prefix, names)
	result, err := awsapi.ListObjects(mockS3API, bucket, prefix)

	assert.Nil(t, err)
	assert.Equal(t, len(result), len(names))

	allNames := [][]string{[]string{"name1", "name2"},
		[]string{"name3", "name4", "name5"}, []string{}}
	truncates := []bool{true, true, false}

	// loop for 3 times to get all results
	mockS3API = mock.S3APIListObjectsPaginated(
		bucket, prefix, allNames, truncates, []string{"", "marker0", "marker1"},
		[]string{"marker0", "marker1", ""})
	result, err = awsapi.ListObjects(mockS3API, bucket, prefix)
	assert.Nil(t, err)
	assert.Equal(t, []string{
		fmt.Sprintf("%s/%s", prefix, "name1"),
		fmt.Sprintf("%s/%s", prefix, "name2"),
		fmt.Sprintf("%s/%s", prefix, "name3"),
		fmt.Sprintf("%s/%s", prefix, "name4"),
		fmt.Sprintf("%s/%s", prefix, "name5")}, result)

	mockS3API = mock.S3APIListObjectsError()
	_, err = awsapi.ListObjects(mockS3API, bucket, prefix)

	assert.NotNil(t, err)

}

func TestListObjectsPaginated(t *testing.T) {

	bucket := "example-bucket"
	prefix := "example-prefix/prefix/example"
	awsapi := AWSAPI{}

	names := [][]string{[]string{"name1", "name2"},
		[]string{"name3", "name4"}, []string{}}
	truncates := []bool{true, true, true}

	mockS3API := mock.S3APIListObjectsPaginated(
		bucket, prefix, names, truncates, []string{"marker", "", "fakeNextMarker"},
		[]string{"", "fakeNextMarker", ""})
	contents, nextMarker, truncated, err := awsapi.ListObjectsPaginated(
		mockS3API, bucket, prefix, 2, "")
	assert.Nil(t, err)
	assert.Equal(t, 2, len(contents))
	assert.Equal(t, true, truncated)
	assert.Equal(t, "", nextMarker)

	contents, nextMarker, truncated, err = awsapi.ListObjectsPaginated(
		mockS3API, bucket, prefix, 2, nextMarker)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(contents))
	assert.Equal(t, true, truncated)
	assert.Equal(t, "fakeNextMarker", nextMarker)

	contents, nextMarker, truncated, err = awsapi.ListObjectsPaginated(
		mockS3API, bucket, prefix, 2, nextMarker)
	assert.Nil(t, err)
	assert.Equal(t, true, truncated)
	assert.Equal(t, 0, len(contents))
	assert.Equal(t, "", nextMarker)

	mockS3API = mock.S3APIListObjectsError()
	contents, nextMarker, truncated, err = awsapi.ListObjectsPaginated(
		mockS3API, bucket, prefix, 2, nextMarker)
	assert.Equal(t, 0, len(contents))
	assert.NotNil(t, err)
	assert.Equal(t, false, truncated)
	assert.Equal(t, "", nextMarker)
	assert.Equal(t, "FakeListObjectsError", err.Error())
}

func TestDeleteObject(t *testing.T) {
	bucket := "example-bucket"
	key := "example-key/example.json"
	awsapi := AWSAPI{}

	versionID := "ESGHby71m6xT4QoerT81suYiCzOod8v_"
	mockS3API := mock.S3APIDeleteObject(versionID)
	err := awsapi.DeleteObject(mockS3API, bucket, key)
	assert.Nil(t, err)
}

func TestDeleteObjectError(t *testing.T) {
	bucket := "example-bucket"
	key := "example-key/example.json"
	awsapi := AWSAPI{}

	mockS3API := mock.S3APIDeleteObjectError()
	err := awsapi.DeleteObject(mockS3API, bucket, key)
	assert.NotNil(t, err)
	assert.Equal(t, "FakeDeleteObjectError", err.Error())
}

func TestPutObject(t *testing.T) {
	bucket := "example-bucket"
	key := "example-key/example.json"
	content := "\"single_channel_guests\":[]"

	versionID := "ESGHby71m6xT4QoerT81suYiCzOod8v_"
	awsapi := AWSAPI{}
	uploader := mock.S3ManagerAPIUpload(versionID)
	err := awsapi.PutObject(uploader, bucket, key, content)
	assert.Nil(t, err)
}

func TestPutObjectError(t *testing.T) {
	bucket := "example-bucket"
	key := "example-key/example.json"
	content := "\"single_channel_guests\":[]"

	awsapi := AWSAPI{}
	uploader := mock.S3ManagerAPIUploadError()
	err := awsapi.PutObject(uploader, bucket, key, content)
	assert.NotNil(t, err)
	assert.Equal(t, "FakeUploadError", err.Error())
}
