package utils

import (
	"bytes"
	"errors"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/xinnige/asteraceae/calendula/mock"
)

func TestReadFile(t *testing.T) {
	_, err := ReadFile("../../README.md")
	assert.Nil(t, err)
}

func TestReadFileError(t *testing.T) {
	noSuchFile := "no_such_file"
	_, err := ReadFile(noSuchFile)
	assert.NotNil(t, err)
}

func TestMarshal(t *testing.T) {
	testmap := map[string]string{"key01": "testValue"} //, "key02": "<tag></tag>&"}
	bytes := Marshal(testmap, &JSONAPI{})
	assert.Equal(t, "{\"key01\":\"testValue\"}", string(bytes))

	bytes = MarshalIndent(testmap, "", "", &JSONAPI{})
	assert.Equal(t, "{\n\"key01\": \"testValue\"\n}\n", string(bytes))
}

func TestMarshalError(t *testing.T) {
	testmap := ""

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockSIface := mock.NewMockSerialInterface(mockCtrl)
	mockSIface.EXPECT().Marshal(gomock.Any()).Return(
		nil, errors.New("FakeJSONError")).Times(1)
	mockSIface.EXPECT().MarshalIndent(
		gomock.Any(), gomock.Any(), gomock.Any()).Return(
		nil, errors.New("FakeJSONError")).Times(1)

	bytes := Marshal(testmap, mockSIface)
	assert.Nil(t, bytes)

	bytes = MarshalIndent(testmap, "", "", mockSIface)
	assert.Nil(t, bytes)
}

func TestUnmarshalJSON(t *testing.T) {
	data := []byte("{\"key01\":\"testValue\"}")
	testmap := map[string]string{}
	err := Unmarshal(data, &testmap, &JSONAPI{})
	assert.Equal(t, "testValue", testmap["key01"])
	assert.Nil(t, err)
}

func TestUnmarshalError(t *testing.T) {
	data := []byte{}
	testmap := map[string]string{}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockSIface := mock.NewMockSerialInterface(mockCtrl)
	mockSIface.EXPECT().Unmarshal(gomock.Any(), gomock.Any()).Return(
		errors.New("FakeJSONError")).Times(1)

	err := Unmarshal(data, &testmap, mockSIface)
	assert.NotNil(t, err)
	assert.Equal(t, "FakeJSONError", err.Error())
}

func TestRegexpMatch(t *testing.T) {
	var match bool
	TestData := [][]string{
		{"invite all", "<@UCJ5LC1M0> invite all members"},
		{"全員招待", "<@UCJ5LC1M0> 全員招待して"},
		{"invite all", "<@UCJ5LC1M0> invite"},
		{"全員招待", "<@UCJ5LC1M0> 招待して"},
		{"/static/", "http://127.0.0.1:8080/static/css/style.css"},
	}
	TestResult := []bool{true, true, false, false, true}

	for i := range TestData {
		match = RegexpIsAnyMatch(TestData[i][0], []byte(TestData[i][1]))
		assert.Equal(t, TestResult[i], match)
	}
}

func TestRegexpMatchError(t *testing.T) {
	match := RegexpIsAnyMatch(`\`, []byte("123"))
	assert.Equal(t, false, match)
}

func TestRegexpSubMatch(t *testing.T) {
	ret := RegexpSubMatch(`<#([A-Z0-9]+)`, []byte("<#ACDESFJS"))
	assert.Equal(t, "ACDESFJS", string(ret))

	ret = RegexpSubMatch(`^arn:aws:([a-z0-9]+)`, []byte(
		"arn:aws:lambda:fake-func"))
	assert.Equal(t, "lambda", string(ret))

}

func TestRegexpSubMatchError(t *testing.T) {
	ret := RegexpSubMatch(`\`, []byte("123"))
	assert.Equal(t, 0, len(ret))
}

func TestRegexpSubMatchNoMatch(t *testing.T) {
	ret := RegexpSubMatch(`[a-z]+`, []byte("123"))
	assert.Equal(t, 0, len(ret))
}

func TestIsStringItemInArray(t *testing.T) {
	item := "Item000"
	arrayIn := []string{item, "Item001", "Item002", "Item003"}
	arrayNotIn := []string{"Item001", "Item002", "Item003"}
	assert.Equal(t, true, IsStringItemInArray(item, arrayIn))
	assert.Equal(t, false, IsStringItemInArray(item, arrayNotIn))
}

func TestMergeArrays(t *testing.T) {
	array1 := []string{"Item000", "Item001", "Item002", "Item003"}
	array2 := []string{"Item001", "Item002", "Item004", "Item005"}
	ret := MergeArrays(array1, array2)
	assert.Equal(t, 6, len(ret))
}

func TestSplitPath(t *testing.T) {
	base, file := SplitPath("/path/path/parent/filename")
	assert.Equal(t, "parent", base)
	assert.Equal(t, "filename", file)

	base, file = SplitPath("/")
	assert.Equal(t, "/", base)
	assert.Equal(t, "", file)
}

func TestGetEnv(t *testing.T) {
	envKey := "eKey"
	envValue := "eValue"
	deValue := "defaultValue"
	os.Setenv(envKey, "")
	assert.Equal(t, deValue, GetEnv(envKey, deValue))
	os.Setenv(envKey, envValue)
	assert.Equal(t, envValue, GetEnv(envKey, deValue))
}

func TestGetEnvInt(t *testing.T) {
	envKey := "eKey01"
	envValue := "10"
	envValueError := "anc"
	deValue := 1
	os.Setenv(envKey, "")
	assert.Equal(t, deValue, GetEnvInt(envKey, deValue))
	os.Setenv(envKey, envValue)
	assert.Equal(t, 10, GetEnvInt(envKey, deValue))
	os.Setenv(envKey, envValueError)
	assert.Equal(t, deValue, GetEnvInt(envKey, deValue))
}

func TestReadFrom(t *testing.T) {
	var buffer bytes.Buffer
	expected := "contents of body"
	reader := bytes.NewBufferString(expected)
	ret, err := ReadFrom(&buffer, reader)
	assert.Nil(t, err)
	assert.Equal(t, expected, ret)
}

func TestReadFromError(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockBuffer := mock.NewMockBufferInterface(mockCtrl)
	mockBuffer.EXPECT().ReadFrom(gomock.Any()).Return(
		int64(0), errors.New("FakeBufferReadError")).Times(1)

	reader := bytes.NewBufferString("")
	ret, err := ReadFrom(mockBuffer, reader)
	assert.Equal(t, "", ret)
	assert.NotNil(t, err)
	assert.Equal(t, "FakeBufferReadError", err.Error())
}

func TestEncryptTextAES(t *testing.T) {
	encoded := "DfdjIjQDhx9Nh5aWA13fuQQsIGsk8yYUWrm4dbPwokA="
	plainText := []byte("plain text with several words")
	cypherText := EncryptTextAES(DecodeBase64(encoded), plainText, &CommonAPI{})

	assert.Equal(t, 45, len(cypherText))
}

func TestEncryptTextAESEncoded(t *testing.T) {
	encoded := "DfdjIjQDhx9Nh5aWA13fuQQsIGsk8yYUWrm4dbPwokA="
	plainText := []byte("plain text with several words")
	cypherText := EncryptTextAESEncoded(encoded, plainText, &CommonAPI{})
	assert.Equal(t, 45, len(DecodeBase64(cypherText)))
}

func TestDecodeBase64(t *testing.T) {
	plain := "plain text with several words"
	encoded := "cGxhaW4gdGV4dCB3aXRoIHNldmVyYWwgd29yZHM="
	decoded := DecodeBase64(encoded)
	assert.Equal(t, plain, string(decoded))
}

func TestDecryptTextAES(t *testing.T) {
	encoded := "DfdjIjQDhx9Nh5aWA13fuQQsIGsk8yYUWrm4dbPwokA="
	cipherText := "Wu5g0oSKIU/o91EClgB8qSRV2XlVdCGHqPKWwPuXDumoXvqO2unVKWo6t10S"
	plainText := DecryptTextAES(DecodeBase64(encoded), DecodeBase64(cipherText))

	expected := []byte("plain text with several words")
	assert.Equal(t, expected, plainText)

	plainText = DecryptTextAES(DecodeBase64(encoded), []byte{})
	assert.Equal(t, []byte{}, plainText)
}

func TestDecryptTextAESEncode(t *testing.T) {
	encoded := "DfdjIjQDhx9Nh5aWA13fuQQsIGsk8yYUWrm4dbPwokA="
	cipherText := "Wu5g0oSKIU/o91EClgB8qSRV2XlVdCGHqPKWwPuXDumoXvqO2unVKWo6t10S"
	plainText := DecryptTextAESEncoded(encoded, cipherText)

	expected := "plain text with several words"
	assert.Equal(t, expected, plainText)
}

func TestDecodeBase64Error(t *testing.T) {
	result := DecodeBase64("XXXXX")
	assert.Nil(t, result)
}

func TestDecryptTextAESError(t *testing.T) {
	corruptedKey := []byte("00000")
	plainText := []byte("")
	ret := DecryptTextAES(corruptedKey, plainText)
	assert.Equal(t, 0, len(ret))
}

func TestDecryptTextAESEncodeError(t *testing.T) {
	corruptedKey := "00000"
	plainText := ""
	ret := DecryptTextAESEncoded(corruptedKey, plainText)
	assert.Equal(t, "", ret)
}

func TestEncryptTextAESError(t *testing.T) {
	corruptedKey := []byte("00000")
	plainText := []byte("")
	ret := EncryptTextAES(corruptedKey, plainText, &CommonAPI{})
	assert.Nil(t, ret)

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockIOIface := mock.NewMockIOInterface(mockCtrl)
	mockIOIface.EXPECT().ReadFull(gomock.Any(), gomock.Any()).Return(
		0, errors.New("FakeIOError")).Times(1)

	key := []byte("1234567890123456")
	ret = EncryptTextAES(key, plainText, mockIOIface)
	assert.Nil(t, ret)
}

func TestEncryptTextAESEncodedError(t *testing.T) {
	corruptedKey := "xxxxxxxxxx9Nh5aWA13fuQQsIGsk8yYUWrm4dbPxxxxx"
	plainText := []byte("")
	ret := EncryptTextAESEncoded(corruptedKey, plainText, &CommonAPI{})
	assert.Equal(t, "", ret)

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockIOIface := mock.NewMockIOInterface(mockCtrl)
	mockIOIface.EXPECT().ReadFull(gomock.Any(), gomock.Any()).Return(
		0, errors.New("FakeIOError")).Times(1)

	key := "DfdjIjQDhx9Nh5aWA13fuQQsIGsk8yYUWrm4dbPwokA="
	ret = EncryptTextAESEncoded(key, plainText, mockIOIface)
	assert.Equal(t, "", ret)

}

func TestTrimExtension(t *testing.T) {
	assert.Equal(t, "2018-09-01", TrimExtension("fake/prefix/2018-09-01.json"))
	assert.Equal(t, "2018-09-01", TrimExtension("2018-09-01.json"))
	assert.Equal(t, "2018-09-01", TrimExtension("2018-09-01."))
	assert.Equal(t, "2018-09-01", TrimExtension("2018-09-01"))
	assert.Equal(t, "", TrimExtension("."))
	assert.Equal(t, "", TrimExtension(""))
}

func TestGetOrderedMapKeys(t *testing.T) {
	originMap := map[string]interface{}{
		"あ":   "ひらがな",
		"99":  "value99",
		"平":   "あ",
		"102": "value102",
		"bbb": "valueBBB",
		"aaa": "valueAAA",
	}
	keys := GetOrderedMapKeys(originMap)
	assert.Equal(t, []string{"102", "99", "aaa", "bbb", "あ", "平"}, keys)
}

func TestRandID(t *testing.T) {
	id := RandID()
	assert.Equal(t, 15, len(id))
}

func TestShortID(t *testing.T) {
	id := ShortID()
	assert.Equal(t, 10, len(id))
}

func TestChunkStringArray(t *testing.T) {
	src := []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"}
	divided := ChunkStringArray(src, 3)
	assert.Equal(t, [][]string{[]string{"1", "2", "3"}, []string{"4", "5", "6"},
		[]string{"7", "8", "9"}, []string{"10"}}, divided)
}
