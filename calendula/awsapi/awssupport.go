package awsapi

import (
	"log"

	"github.com/aws/aws-sdk-go/service/support"
	"github.com/aws/aws-sdk-go/service/support/supportiface"
)

// NewSupportAPI returns a s3 client
func NewSupportAPI(sess AWSSession) supportiface.SupportAPI {
	awssess, err := sess.NewSession()
	if err != nil {
		log.Printf("AWSError: Cannot create aws session. err %#v", err)
		return nil
	}
	return support.New(awssess)
}
