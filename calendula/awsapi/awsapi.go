package awsapi

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kms/kmsiface"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/aws/aws-sdk-go/service/s3/s3manager/s3manageriface"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"
)

// AWSSession defines an interface of session.NewSession
type AWSSession interface {
	NewSession() (*session.Session, error)
}

// AWSServiceSession defiens a struct to implement AWSSession
type AWSServiceSession struct {
}

// NewSession implements AWSSession.NewSession
func (awsSession *AWSServiceSession) NewSession() (*session.Session, error) {
	return session.NewSession()
}

// AWSAPI defines a struct to implement AWSInterface
type AWSAPI struct {
}

// AWSIface holds aws clients
type AWSIface struct {
	KMSSVC       kmsiface.KMSAPI
	S3ManagerSVC s3manageriface.UploaderAPI
	S3SVC        s3iface.S3API
	SSMSVC       ssmiface.SSMAPI
}
