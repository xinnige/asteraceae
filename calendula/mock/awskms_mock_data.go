package mock

import (
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/aws/aws-sdk-go/service/kms/kmsiface"
)

// KMSAPI mocks s3iface.S3API
type KMSAPI struct {
	kmsiface.KMSAPI
	genDataKeyOutput kms.GenerateDataKeyOutput
	genErr           error
	decryptOutput    kms.DecryptOutput
	decryptErr       error
}

// GenerateDataKey mocks KMSAPI.GenerateDataKey
func (m KMSAPI) GenerateDataKey(
	*kms.GenerateDataKeyInput) (*kms.GenerateDataKeyOutput, error) {
	return &m.genDataKeyOutput, m.genErr
}

// Decrypt mocks KMSAPI.Decrypt
func (m KMSAPI) Decrypt(*kms.DecryptInput) (*kms.DecryptOutput, error) {
	return &m.decryptOutput, m.decryptErr
}

func mockDecryptOutput(plain []byte, keyID string) *kms.DecryptOutput {
	return &kms.DecryptOutput{
		Plaintext: plain,
		KeyId:     aws.String(keyID),
	}
}

// KMSAPIDecrypt returns KMSAPI with KMSAPI.Decrypt
func KMSAPIDecrypt(plain []byte, keyID string) *KMSAPI {
	return &KMSAPI{
		decryptOutput: *mockDecryptOutput(plain, keyID),
		decryptErr:    nil,
	}
}

// KMSAPIDecryptError returns KMSAPI with KMSAPI.Decrypt error
func KMSAPIDecryptError() *KMSAPI {
	return &KMSAPI{
		decryptOutput: kms.DecryptOutput{},
		decryptErr:    errors.New("FakeKMSDecryptError"),
	}
}

// FakeDataKey returns a data key
func FakeDataKey() []byte {
	return []byte{177, 78, 243, 127, 250, 70, 35, 198, 136, 145, 116, 110,
		201, 95, 121, 38, 50, 83, 78, 210, 135, 118, 6, 81, 196, 72, 196, 139,
		55, 209, 55, 183}
}

func mockGenDataKeyOutput(
	cipher []byte, plain []byte, keyID string) *kms.GenerateDataKeyOutput {
	return &kms.GenerateDataKeyOutput{
		CiphertextBlob: cipher,
		Plaintext:      plain,
		KeyId:          aws.String(keyID),
	}
}

// KMSAPIGenDataKey returns mock.KMSAPI with KMSAPI.GenerateDataKey
func KMSAPIGenDataKey(cipher []byte, plain []byte, keyID string) *KMSAPI {
	return &KMSAPI{
		genDataKeyOutput: *mockGenDataKeyOutput(cipher, plain, keyID),
		genErr:           nil,
	}
}

// KMSAPIGenDataKeyError returns KMSAPI with KMSAPI.GenerateDataKey error
func KMSAPIGenDataKeyError() *KMSAPI {
	return &KMSAPI{
		genDataKeyOutput: kms.GenerateDataKeyOutput{},
		genErr:           errors.New("FakeGenDataKeyError"),
	}
}
