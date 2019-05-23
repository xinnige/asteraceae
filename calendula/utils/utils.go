package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"

	"github.com/rs/xid"
	"github.com/teris-io/shortid"
)

// IOInterface defines an interface for io operations
type IOInterface interface {
	ReadFull(r io.Reader, buf []byte) (n int, err error)
}

// BufferInterface defines an interface for bytes.Buffer
type BufferInterface interface {
	ReadFrom(r io.Reader) (n int64, err error)
	String() string
}

// SerialInterface defines an interface for serialization
type SerialInterface interface {
	Marshal(v interface{}) ([]byte, error)
	MarshalIndent(v interface{}, prefix, indent string) ([]byte, error)
	Unmarshal(data []byte, v interface{}) error
}

// ReadCloser mocks io.ReadCloser
type ReadCloser struct {
	io.Reader
}

// Close mocks io.Closer
func (ReadCloser) Close() error { return nil }

// CommonAPI defines a struct for io IOInterface
type CommonAPI struct {
}

// ReadFull wraps io.ReadFull
func (api *CommonAPI) ReadFull(r io.Reader, buf []byte) (n int, err error) {
	return io.ReadFull(r, buf)
}

// JSONAPI defines a struct for mock yaml
type JSONAPI struct {
}

// Marshal wraps json.Marshal
func (api *JSONAPI) Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

// Unmarshal wraps json.Unmarshal
func (api *JSONAPI) Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

// MarshalIndent wraps json.MarshalIndent
func (api *JSONAPI) MarshalIndent(
	v interface{}, prefix, indent string) ([]byte, error) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	err := enc.Encode(v)
	if err != nil {
		return nil, err
	}

	var inbuf bytes.Buffer
	err = json.Indent(&inbuf, buf.Bytes(), prefix, indent)
	return inbuf.Bytes(), err
	// return json.MarshalIndent(v, prefix, indent)
}

// ReadFile helps to read all lines in a file
func ReadFile(filename string) ([]byte, error) {
	content, err := ioutil.ReadFile(filepath.Clean(filename))
	if err != nil {
		log.Printf("Error: ioutils.ReadFile error %#v", err)
		return nil, err
	}
	return content, nil
}

// Marshal helps return a marshalled byte array
func Marshal(v interface{}, siface SerialInterface) []byte {
	bytes, err := siface.Marshal(v)
	if err != nil {
		log.Printf("Error: %T Marshal(%+v) error %#v", siface, v, err)
		return nil
	}
	return bytes
}

// Unmarshal helps unmarshal a json-format []byte
func Unmarshal(data []byte, v interface{}, siface SerialInterface) error {
	err := siface.Unmarshal(data, v)
	if err != nil {
		log.Printf("Error: %T Unmarshal(%v) error %#v", siface, v, err)
		return err
	}
	return nil
}

// MarshalIndent helps return a indented marshalled byte array
func MarshalIndent(
	v interface{}, prefix string, indent string, siface SerialInterface) []byte {
	bytes, err := siface.MarshalIndent(v, prefix, indent)
	if err != nil {
		log.Printf("Error: %T Marshal(%+v) error %#v", siface, v, err)
		return nil
	}
	return bytes
}

// RegexpIsAnyMatch helps to match substring and returns a boolean value
func RegexpIsAnyMatch(regex string, text []byte) bool {
	re, err := regexp.Compile(regex)
	if err != nil {
		log.Printf("Error: cannot match regexp %s in %s, err %#v.\n",
			regex, text, err)
		return false
	}
	return re.Match(text)
}

// RegexpSubMatch helps to get the second match in a byte array
func RegexpSubMatch(regex string, text []byte) []byte {
	re, err := regexp.Compile(regex)
	if err != nil {
		log.Printf(
			"Error: cannot find any match of regular expression %s in %s, err %v",
			regex, text, err)
		return []byte{}
	}
	matches := re.FindSubmatch(text)
	if len(matches) < 2 {
		return []byte{}
	}
	return matches[1]
}

// IsStringItemInArray checks if item in array or not
func IsStringItemInArray(item string, array []string) bool {
	for index := range array {
		if item == array[index] {
			return true
		}
	}
	return false
}

// MergeArrays helps to merge array2 to array1 w/ no duplication
func MergeArrays(array1 []string, array2 []string) []string {
	for idx := range array2 {
		if len(array2[idx]) == 0 {
			continue
		}
		if IsStringItemInArray(array2[idx], array1) {
			continue
		}
		array1 = append(array1, array2[idx])
	}
	return array1
}

// SplitPath helps split filepath to a pair of parent dir and file
func SplitPath(path string) (string, string) {
	dir, file := filepath.Split(path)
	base := filepath.Base(dir)
	return base, file
}

// GetEnv returns an env variable as string or a default value if no env found
func GetEnv(envKey string, deValue string) string {
	value := os.Getenv(envKey)
	if len(value) == 0 {
		return deValue
	}
	return value
}

// GetEnvInt returns an env variable as int or a default value if no env found
func GetEnvInt(envKey string, deValue int) int {
	value := os.Getenv(envKey)
	if len(value) == 0 {
		return deValue
	}
	v, err := strconv.Atoi(value)
	if err != nil {
		log.Printf("Error: cannot convert string %s to int.\n", value)
		return deValue
	}
	return v
}

// ReadFrom returns a string read from io.Reader
func ReadFrom(bufiface BufferInterface, r io.Reader) (string, error) {
	_, err := bufiface.ReadFrom(r)
	if err != nil {
		log.Printf("Error: cannot read from %v, %+v\n", r, err)
		return "", err
	}

	return bufiface.String(), nil
}

// EncodeBase64 returns 64-base encoded text
func EncodeBase64(decoded []byte) string {
	return base64.StdEncoding.EncodeToString(decoded)
}

// DecodeBase64 returns a byte array after decode 64-base text
func DecodeBase64(encoded string) []byte {
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		log.Printf("Error: cannot decode base64 text %s, err: %+v\n", encoded, err)
		return nil
	}
	return decoded
}

// EncryptTextAES encrypts plain text with AES
func EncryptTextAES(key []byte, plainText []byte, ioiface IOInterface) []byte {
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Printf("Error: cannot encrypt AES, err: %+v\n", err)
		return nil
	}

	cipherText := make([]byte, aes.BlockSize+len(plainText))
	iv := cipherText[:aes.BlockSize]
	if _, err := ioiface.ReadFull(rand.Reader, iv); err != nil {
		log.Printf("Error: cannot encrypt AES, err: %+v\n", err)
		return nil
	}

	encryptStream := cipher.NewCTR(block, iv)
	encryptStream.XORKeyStream(cipherText[aes.BlockSize:], plainText)
	// log.Printf("AES-encrypted cipher text length: %d \n", len(cipherText))

	return cipherText
}

// DecryptTextAES decrypts cipher text with AES
func DecryptTextAES(key []byte, cipherText []byte) []byte {
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Printf("Error: cannot decrypt AES, err: %+v\n", err)
		return []byte{}
	}

	if aes.BlockSize >= len(cipherText) {
		log.Printf("Error: cannot decrypt AES, invalid cipher length %d (<= %d)\n",
			len(cipherText), aes.BlockSize)
		return []byte{}
	}
	decryptedText := make([]byte, len(cipherText[aes.BlockSize:]))
	decryptStream := cipher.NewCTR(block, cipherText[:aes.BlockSize])
	decryptStream.XORKeyStream(decryptedText, cipherText[aes.BlockSize:])
	return decryptedText
}

// EncryptTextAESEncoded encrypts plain text with base64 AES key
func EncryptTextAESEncoded(
	encodeKey string, plainText []byte, ioiface IOInterface) string {
	plainKey := DecodeBase64(encodeKey)
	if plainKey == nil {
		return ""
	}
	cypherText := EncryptTextAES(plainKey, plainText, ioiface)
	if cypherText == nil {
		return ""
	}
	return EncodeBase64(cypherText)
}

// DecryptTextAESEncoded decrypts cipher w/ base64 AES EncryptTextAES
func DecryptTextAESEncoded(encodeKey string, cipherText string) string {
	plainKey := DecodeBase64(encodeKey)
	decodeCipherText := DecodeBase64(cipherText)
	if len(plainKey) == 0 || len(decodeCipherText) == 0 {
		return ""
	}
	plainText := DecryptTextAES(plainKey, decodeCipherText)
	return string(plainText)
}

// TrimExtension removes extension of a file name
func TrimExtension(path string) string {
	_, filename := filepath.Split(path)
	ext := filepath.Ext(filename)
	return filename[0 : len(filename)-len(ext)]
}

// GetOrderedMapKeys returns an list of ordered keys of a map
func GetOrderedMapKeys(originMap map[string]interface{}) []string {
	keys := make([]string, len(originMap))
	idx := 0
	for key := range originMap {
		keys[idx] = key
		idx++
	}
	sort.Strings(keys)
	return keys
}

// RandID returns an auto-generated id
func RandID() string {
	uid := xid.New().String()
	if sid, err := shortid.Generate(); err == nil {
		return uid[14:] + sid[:9]
	}
	return uid
}

// ShortID returns an auto-generated short-id
func ShortID() string {
	if sid, err := shortid.Generate(); err == nil {
		return sid
	}
	uid := xid.New().String()
	return uid[10:]
}

// ChunkStringArray by size
func ChunkStringArray(srcArray []string, chunkSize int) [][]string {
	var divided [][]string
	srcLen := len(srcArray)
	for i := 0; i < srcLen; i += chunkSize {
		end := i + chunkSize
		if end > srcLen {
			end = srcLen
		}
		divided = append(divided, srcArray[i:end])
	}
	return divided
}
