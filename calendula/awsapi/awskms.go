package awsapi

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/aws/aws-sdk-go/service/kms/kmsiface"
)

// NewKMSAPI returns a KMS API
func NewKMSAPI(sess AWSSession) kmsiface.KMSAPI {
	ksess, err := sess.NewSession()
	if err != nil {
		log.Printf("AWSError: Cannot create aws session, err %#v", err)
		return nil
	}
	return kms.New(ksess)
}

// NewDataKey returns plain datakey with ciphertext-blob
func (awsapi *AWSAPI) NewDataKey(svc kmsiface.KMSAPI, keyAlias string, keySpec string) ([]byte, []byte, error) {
	input := &kms.GenerateDataKeyInput{
		KeyId:   aws.String(keyAlias),
		KeySpec: aws.String(keySpec),
	}
	result, err := svc.GenerateDataKey(input)
	if err != nil {
		log.Printf("AWSError: Cannot generate a new data key, err %+v\n", err)
		return nil, nil, err
	}
	log.Printf("Generated a datakey, length %d", len(result.Plaintext))
	return result.Plaintext, result.CiphertextBlob, nil
}

// GetDataKey returns plain datakey
func (awsapi *AWSAPI) GetDataKey(svc kmsiface.KMSAPI, ciphertext []byte) ([]byte, error) {
	input := &kms.DecryptInput{
		CiphertextBlob: ciphertext,
	}
	result, err := svc.Decrypt(input)
	if err != nil {
		log.Printf("AWSError: Cannot decrypt to get data key, err %+v\n", err)
		return nil, err
	}
	log.Printf("Get the decrypted data key, length %d\n", len(result.Plaintext))
	return result.Plaintext, nil
}
