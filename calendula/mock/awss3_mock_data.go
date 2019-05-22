package mock

import (
	"bytes"
	"errors"
	"fmt"
	io "io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/aws/aws-sdk-go/service/s3/s3manager/s3manageriface"
)

// S3API mocks s3iface.S3API
type S3API struct {
	s3iface.S3API
	getIndex           *int
	getObjectOutputs   []s3.GetObjectOutput
	getErrors          []error
	listIndex          *int
	listObjectsOutputs []s3.ListObjectsV2Output
	listErrors         []error
	deleteIndex        *int
	deleteObjectOutput *s3.DeleteObjectOutput
	deleteErrors       []error
}

// S3ManagerAPI mocks s3manageriface.UploaderAPI
type S3ManagerAPI struct {
	s3manageriface.UploaderAPI
	uploadOutput s3manager.UploadOutput
	err          error
}

// GetObject mocks S3API.GetObject
func (m S3API) GetObject(*s3.GetObjectInput) (*s3.GetObjectOutput, error) {
	var output s3.GetObjectOutput
	if *m.getIndex < len(m.getObjectOutputs) {
		output = m.getObjectOutputs[*m.getIndex]
		*m.getIndex = (*m.getIndex + 1) % len(m.getObjectOutputs)
	}
	var err error
	if *m.getIndex < len(m.getErrors) {
		err = m.getErrors[*m.getIndex]
	}
	return &output, err
}

// ListObjectsV2 mocks S3API.ListObjectsV2
func (m S3API) ListObjectsV2(
	*s3.ListObjectsV2Input) (*s3.ListObjectsV2Output, error) {
	var output s3.ListObjectsV2Output
	if *m.listIndex < len(m.listObjectsOutputs) {
		output = m.listObjectsOutputs[*m.listIndex]
		*m.listIndex = (*m.listIndex + 1) % len(m.listObjectsOutputs)
	}
	var err error
	if *m.listIndex < len(m.listErrors) {
		err = m.listErrors[*m.listIndex]
	}
	return &output, err
}

// DeleteObject mocks S3API.DeleteObject
func (m S3API) DeleteObject(
	*s3.DeleteObjectInput) (*s3.DeleteObjectOutput, error) {
	var err error
	if *m.deleteIndex < len(m.deleteErrors) {
		err = m.deleteErrors[*m.deleteIndex]
		*m.deleteIndex = (*m.deleteIndex + 1) % len(m.deleteErrors)
	}
	return m.deleteObjectOutput, err
}

// Upload mocks UploaderAPI.Upload
func (m S3ManagerAPI) Upload(input *s3manager.UploadInput,
	options ...func(*s3manager.Uploader)) (*s3manager.UploadOutput, error) {
	return &m.uploadOutput, m.err
}

// ReadCloser mocks io.ReadCloser
type ReadCloser struct {
	io.Reader
}

// Close mocks io.Closer
func (ReadCloser) Close() error { return nil }

func mockGetObjectOutput(content string, versionID string) *s3.GetObjectOutput {
	return &s3.GetObjectOutput{
		Body:      ReadCloser{bytes.NewBufferString(content)},
		VersionId: aws.String(versionID),
	}
}

func mockListObjectsOutput(bucket, prefix string, names []string,
	isTruncated bool) *s3.ListObjectsV2Output {
	contents := make([]*s3.Object, len(names))
	for idx, name := range names {
		contents[idx] = &s3.Object{
			Key: aws.String(fmt.Sprintf("%s/%s", prefix, name)),
		}
	}
	return &s3.ListObjectsV2Output{
		Contents:    contents,
		Name:        aws.String(bucket),
		Prefix:      aws.String(prefix),
		IsTruncated: aws.Bool(isTruncated),
	}
}

func mockDeleteObjectOutput(vid string) *s3.DeleteObjectOutput {
	return &s3.DeleteObjectOutput{
		VersionId: aws.String(vid),
	}
}

func mockUploadOutput(versionID string) *s3manager.UploadOutput {
	return &s3manager.UploadOutput{
		VersionID: aws.String(versionID),
	}
}

// S3APIGetObjects returns S3API w/ multiple times of read s3.GetObject
func S3APIGetObjects(contents []string, versionID string) *S3API {
	s3api := &S3API{
		getIndex:         aws.Int(0),
		getObjectOutputs: make([]s3.GetObjectOutput, 0),
	}
	for _, content := range contents {
		s3api.getObjectOutputs = append(
			s3api.getObjectOutputs, *mockGetObjectOutput(content, versionID))
	}
	return s3api
}

// S3APIGetObjectError returns S3API with s3.GetObject
func S3APIGetObjectError() *S3API {
	return &S3API{
		getIndex:         aws.Int(0),
		getObjectOutputs: []s3.GetObjectOutput{},
		getErrors:        []error{errors.New("FakeGetObjectError")},
	}
}

func mockListObjectsOutputPaginated(
	bucket, prefix string, names []string,
	isTruncated bool, marker string, nextMarker string) *s3.ListObjectsV2Output {
	contents := make([]*s3.Object, len(names))
	for idx, name := range names {
		contents[idx] = &s3.Object{
			Key: aws.String(fmt.Sprintf("%s/%s", prefix, name)),
		}
	}
	output := &s3.ListObjectsV2Output{
		Contents:          contents,
		Name:              aws.String(bucket),
		Prefix:            aws.String(prefix),
		ContinuationToken: aws.String(marker),
		IsTruncated:       aws.Bool(isTruncated),
	}
	if len(nextMarker) != 0 {
		output.NextContinuationToken = aws.String(nextMarker)
	}
	return output
}

// S3APIListObjects returns S3API with s3.ListObjects
func S3APIListObjects(bucket, prefix string, names []string) *S3API {
	return &S3API{
		listIndex: aws.Int(0),
		listObjectsOutputs: []s3.ListObjectsV2Output{
			*mockListObjectsOutput(bucket, prefix, names, false)},
	}
}

// S3APIListObjectsPaginated returns S3API with s3.ListObjects truncated
func S3APIListObjectsPaginated(bucket, prefix string, names [][]string,
	truncated []bool, marker []string, next []string) *S3API {
	s3api := &S3API{
		listIndex:          aws.Int(0),
		listObjectsOutputs: make([]s3.ListObjectsV2Output, 0),
	}
	for idx, ns := range names {
		s3api.listObjectsOutputs = append(s3api.listObjectsOutputs,
			*mockListObjectsOutputPaginated(
				bucket, prefix, ns, truncated[idx], marker[idx], next[idx]))
	}
	return s3api
}

// S3APIListObjectsError returns S3API with s3.ListObjects error
func S3APIListObjectsError() *S3API {
	return &S3API{
		listIndex:          aws.Int(0),
		listObjectsOutputs: []s3.ListObjectsV2Output{},
		listErrors:         []error{errors.New("FakeListObjectsError")},
	}
}

// S3APIDeleteObject returns S3API with s3.DeleteObject
func S3APIDeleteObject(versionID string) *S3API {
	return &S3API{
		deleteObjectOutput: mockDeleteObjectOutput(versionID),
		deleteErrors:       nil,
		deleteIndex:        aws.Int(0),
	}
}

// S3APIDeleteObjectError returns S3API with s3.DeleteObject with error
func S3APIDeleteObjectError() *S3API {
	return &S3API{
		deleteObjectOutput: &s3.DeleteObjectOutput{},
		deleteErrors:       []error{errors.New("FakeDeleteObjectError")},
		deleteIndex:        aws.Int(0),
	}
}

// S3APIAll return S3API with s3.ListObjects, s3.GetObject
func S3APIAll(bucket, prefix string, names [][]string,
	contents []string, versionID string, getErr []error, listErr []error) *S3API {
	s3api := &S3API{
		getIndex:    aws.Int(0),
		listIndex:   aws.Int(0),
		deleteIndex: aws.Int(0),
		getErrors:   getErr,
		listErrors:  listErr,
	}
	for _, content := range contents {
		s3api.getObjectOutputs = append(s3api.getObjectOutputs,
			*mockGetObjectOutput(content, versionID))
	}
	for _, ns := range names {
		s3api.listObjectsOutputs = append(s3api.listObjectsOutputs,
			*mockListObjectsOutput(bucket, prefix, ns, false))
	}

	return s3api
}

// S3ManagerAPIUpload returns mock.S3ManagerAPI
func S3ManagerAPIUpload(versionID string) *S3ManagerAPI {
	return &S3ManagerAPI{
		uploadOutput: *mockUploadOutput(versionID),
		err:          nil,
	}
}

// S3ManagerAPIUploadError returns mock.S3ManagerAPI
func S3ManagerAPIUploadError() *S3ManagerAPI {
	return &S3ManagerAPI{
		uploadOutput: s3manager.UploadOutput{},
		err:          errors.New("FakeUploadError"),
	}
}
