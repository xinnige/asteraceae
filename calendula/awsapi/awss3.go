package awsapi

import (
	"bytes"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/aws/aws-sdk-go/service/s3/s3manager/s3manageriface"
	"github.com/xinnige/asteraceae/calendula/utils"
)

const (
	maxsize = 1000
)

// NewS3API returns a s3 client
func NewS3API(sess AWSSession) s3iface.S3API {
	s3ess, err := sess.NewSession()
	if err != nil {
		log.Printf("AWSError: Cannot create aws session. err %#v", err)
		return nil
	}
	return s3.New(s3ess)
}

// NewS3UploaderAPI returns a s3manager.Uploader pointer
func NewS3UploaderAPI(svc s3iface.S3API) s3manageriface.UploaderAPI {
	return s3manager.NewUploaderWithClient(svc)
}

// GetObject returns the string content of a s3 object
func (awsapi *AWSAPI) GetObject(svc s3iface.S3API, bucket string, key string) (string, error) {
	input := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}
	result, err := svc.GetObject(input)
	if err != nil {
		log.Printf("AWSError: Cannot get objects from s3://%s/%s, err: %+v", bucket, key, err)
		return "", err
	}
	log.Printf("Get a object from s3://%s/%s, VersionId %s\n", bucket, key, *result.VersionId)
	var buf bytes.Buffer
	return utils.ReadFrom(&buf, result.Body)
}

func list2array(result *s3.ListObjectsV2Output) []string {
	objects := make([]string, len(result.Contents))
	for idx, obj := range result.Contents {
		objects[idx] = *obj.Key
	}
	return objects
}

// ListObjects helps to list all objects in a bucket with a certain prefix
func (awsapi *AWSAPI) ListObjects(svc s3iface.S3API, bucket string, prefix string) ([]string, error) {
	found := make([]string, 0)
	input := &s3.ListObjectsV2Input{
		Bucket:  aws.String(bucket),
		Prefix:  aws.String(prefix),
		MaxKeys: aws.Int64(maxsize),
	}
	result, err := svc.ListObjectsV2(input)

	// if truncated, loop to get all results
	for err == nil && *result.IsTruncated && result.NextContinuationToken != nil {
		found = append(found, list2array(result)...)
		input = &s3.ListObjectsV2Input{
			Bucket:            aws.String(bucket),
			Prefix:            aws.String(prefix),
			MaxKeys:           aws.Int64(maxsize),
			ContinuationToken: result.NextContinuationToken,
		}
		result, err = svc.ListObjectsV2(input)
	}
	if err != nil {
		log.Printf("AWSError: Cannot list objects from s3://%s/%s, err: %+v", bucket, prefix, err)
		return nil, err
	}
	found = append(found, list2array(result)...)
	return found, nil
}

// ListObjectsPaginated return a truncated list of object keys
func (awsapi *AWSAPI) ListObjectsPaginated(svc s3iface.S3API, bucket string, prefix string, size int, marker string) ([]string, string, bool, error) {
	input := &s3.ListObjectsV2Input{
		Bucket:  aws.String(bucket),
		Prefix:  aws.String(prefix),
		MaxKeys: aws.Int64(int64(size)),
	}
	if len(marker) != 0 {
		input.ContinuationToken = aws.String(marker)
	}
	result, err := svc.ListObjectsV2(input)
	if err != nil {
		log.Printf("AWSError: Cannot list objects from s3://%s/%s, err: %+v", bucket, prefix, err)
		return nil, "", false, err
	}

	var nextMarker string
	if result.NextContinuationToken != nil {
		nextMarker = *result.NextContinuationToken
	}
	return list2array(result), nextMarker, *result.IsTruncated, nil
}

// DeleteObject deletes a object in s3
func (awsapi *AWSAPI) DeleteObject(svc s3iface.S3API, bucket string, key string) error {
	input := &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}
	result, err := svc.DeleteObject(input)
	if err != nil {
		log.Printf("AWSError: Cannot delete object at s3://%s/%s, err: %+v", bucket, key, err)
		return err
	}
	log.Printf("Deleted a object from s3://%s/%s, VersionId %s\n", bucket, key, *result.VersionId)
	return nil
}

// PutObject uploads a string as a s3 object
func (awsapi *AWSAPI) PutObject(uploaderAPI s3manageriface.UploaderAPI, bucket string, key string, content string) error {
	input := &s3manager.UploadInput{
		Body:   strings.NewReader(content),
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}
	result, err := uploaderAPI.Upload(input)
	if err != nil {
		log.Printf("AWSError: Cannot upload object to s3://%s/%s, err: %+v", bucket, key, err)
		return err
	}
	log.Printf("Upload a object to s3://%s/%s, VersionId %s", bucket, key, *result.VersionID)
	return nil
}
