package awsapi

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/xinnige/asteraceae/calendula/mock"
)

func TestNewKMSAPI(t *testing.T) {
	sess := &AWSServiceSession{}
	kmsapi := NewKMSAPI(sess)
	assert.NotNil(t, kmsapi)
}

func TestNewKMSAPIError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockSession := mock.NewMockAWSSession(mockCtrl)
	mockSession.EXPECT().NewSession().Return(
		nil, errors.New("FakeNewSessionError")).Times(1)

	kmsapi := NewKMSAPI(mockSession)
	assert.Nil(t, kmsapi)
}

func TestGetDataKey(t *testing.T) {
	plain := []byte{98, 170, 127, 104, 186, 65, 55, 18, 120, 179, 231, 184, 243, 65, 102, 51, 96, 105, 139, 42, 105, 196, 50, 91, 206, 144, 24, 249, 125, 96, 140, 2}
	keyID := "fake keyid"
	kmsapi := mock.KMSAPIDecrypt(plain, keyID)

	awsapi := AWSAPI{}
	ciphertextBlob := []byte{1, 2, 3, 0, 120, 16, 93, 38, 62, 53, 48, 239, 163, 214, 39, 208, 45, 182, 162, 6, 224, 254, 1, 229, 147, 63, 103, 26, 141, 19, 75, 222, 187, 208, 126, 160, 53, 1, 92, 45, 130, 192, 9, 209, 108, 247, 148, 150, 152, 107, 21, 149, 67, 99, 0, 0, 0, 126, 48, 124, 6, 9, 42, 134, 72, 134, 247, 13, 1, 7, 6, 160, 111, 48, 109, 2, 1, 0, 48, 104, 6, 9, 42, 134, 72, 134, 247, 13, 1, 7, 1, 48, 30, 6, 9, 96, 134, 72, 1, 101, 3, 4, 1, 46, 48, 17, 4, 12, 41, 173, 169, 140, 215, 162, 87, 215, 78, 88, 96, 88, 2, 1, 16, 128, 59, 101, 105, 218, 82, 23, 121, 61, 145, 46, 235, 193, 161, 82, 70, 148, 198, 182, 198, 229, 3, 221, 40, 194, 244, 195, 114, 91, 186, 190, 250, 31, 23, 81, 213, 244, 239, 64, 126, 122, 236, 17, 18, 155, 42, 46, 125, 153, 202, 92, 149, 83, 205, 253, 153, 232, 158, 74, 170, 214}
	plainkey, err := awsapi.GetDataKey(kmsapi, ciphertextBlob)

	assert.Nil(t, err)
	assert.Equal(t, plain, plainkey)
}

func TestGetDataKeyError(t *testing.T) {
	kmsapi := mock.KMSAPIDecryptError()
	awsapi := AWSAPI{}
	ciphertextBlob := []byte{0, 1}
	_, err := awsapi.GetDataKey(kmsapi, ciphertextBlob)

	assert.NotNil(t, err)
	assert.Equal(t, "FakeKMSDecryptError", err.Error())
}

func TestNewDataKey(t *testing.T) {
	cipher := []byte{1, 2, 3, 0, 120, 16, 93}
	plain := []byte{98, 170, 127, 104, 186, 65, 55, 18, 120, 179, 231, 184, 243, 65, 102, 51, 96, 105, 139, 42, 105, 196, 50, 91, 206, 144, 24, 249, 125, 96, 140, 2}
	keyID := "fake keyid"

	kmsapi := mock.KMSAPIGenDataKey(cipher, plain, keyID)

	awsapi := AWSAPI{}
	keyAlias := "alias/slack-admin-tool-encrypter"
	keySpec := "AES_256"
	plainkey, ciphertextBlob, err := awsapi.NewDataKey(kmsapi, keyAlias, keySpec)
	assert.Nil(t, err)
	assert.Equal(t, plain, plainkey)
	assert.Equal(t, cipher, ciphertextBlob)

}

func TestNewDataKeyError(t *testing.T) {
	kmsapi := mock.KMSAPIGenDataKeyError()
	awsapi := AWSAPI{}
	keyAlias := "alias/slack-admin-tool-encrypter"
	keySpec := "AES_256"
	_, _, err := awsapi.NewDataKey(kmsapi, keyAlias, keySpec)
	assert.NotNil(t, err)
	assert.Equal(t, "FakeGenDataKeyError", err.Error())
}
