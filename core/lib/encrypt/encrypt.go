package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/md5" // #nosec
	"crypto/rand"
	"crypto/sha1" // #nosec
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"hash"
	"hash/crc32"
	"io"
)

var (
	AESKEY       = "welcome mshk top 32bytes32bytes!"
	JWTSecretKey = []byte("welcome mshk top")
)

// JWT 加密
func JWTEncrypt(mapClaims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, mapClaims)
	return token.SignedString(JWTSecretKey)
}

// JWT 解密
func JWTDecrypt(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return JWTSecretKey, nil
	})
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}

// AES 加密
func AesEncrypt(orig, key string) (string, error) {
	// generate a new aes cipher using our 32 byte long key
	c, err := aes.NewCipher([]byte(key))
	// if there are any errors, handle them
	if err != nil {
		return "", err
	}

	// gcm or Galois/Counter Mode, is a mode of operation
	// for symmetric key cryptographic block ciphers
	// - https://en.wikipedia.org/wiki/Galois/Counter_Mode
	gcm, err := cipher.NewGCM(c)
	// if any error generating new GCM
	// handle them
	if err != nil {
		return "", err
	}

	// creates a new byte array the size of the nonce
	// which must be passed to Seal
	nonce := make([]byte, gcm.NonceSize())
	// populates our nonce with a cryptographically secure
	// random sequence
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	// here we encrypt our text using the Seal function
	// Seal encrypts and authenticates plaintext, authenticates the
	// additional data and appends the result to dst, returning the updated
	// slice. The nonce must be NonceSize() bytes long and unique for all
	// time, for a given key.
	// signatured := commonutils.HexEncodeToString()
	return Base64Encode(gcm.Seal(nonce, nonce, []byte(orig), nil)), nil

}

// AES 解密
func AesDecrypt(cryted, key string) (string, error) {
	ciphertext, err := Base64Decode(cryted)
	// if our program was unable to read the file
	// print out the reason why it can't
	if err != nil {
		return "", err
	}

	c, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", err
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}
	return string(plaintext), nil
}

// Const declarations for common.go operations
const (
	HashSHA1 = iota
	HashSHA256
	HashSHA512
	HashSHA512_384
	MD5New
)

//	数字版本的md5
func CRC32(str string) uint32 {
	return crc32.ChecksumIEEE([]byte(str))
}

// GetRandomSalt returns a random salt
func GetRandomSalt(input []byte, saltLen int) ([]byte, error) {
	if saltLen <= 0 {
		return nil, errors.New("salt length is too small")
	}
	salt := make([]byte, saltLen)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return nil, err
	}

	var result []byte
	if input != nil {
		result = input
	}
	result = append(result, salt...)
	return result, nil
}

// GetMD5 returns a MD5 hash of a byte array
func GetMD5(input []byte) ([]byte, error) {
	shaobj := md5.New() // #nosec
	if _, err := shaobj.Write(input); err != nil {
		return nil, err
	}
	return shaobj.Sum(nil), nil
}

// GetSHA512 returns a SHA512 hash of a byte array
func GetSHA512(input []byte) ([]byte, error) {
	shaobj := sha512.New()
	if _, err := shaobj.Write(input); err != nil {
		return nil, err
	}
	return shaobj.Sum(nil), nil
}

// GetSHA256 returns a SHA256 hash of a byte array
func GetSHA256(input []byte) ([]byte, error) {
	shaobj := sha256.New()
	if _, err := shaobj.Write(input); err != nil {
		return nil, err
	}
	return shaobj.Sum(nil), nil
}

// GetHMAC returns a keyed-hash message authentication code using the desired
// hashtype
func GetHMAC(hashType int, input, key []byte) ([]byte, error) {
	var hashFunc func() hash.Hash

	switch hashType {
	case HashSHA1:
		hashFunc = sha1.New
	case HashSHA256:
		hashFunc = sha256.New
	case HashSHA512:
		hashFunc = sha512.New
	case HashSHA512_384:
		hashFunc = sha512.New384
	case MD5New:
		hashFunc = md5.New
	}

	// 使用给定的hash.Hash类型和密钥返回新的HMAC哈希
	hmacFunc := hmac.New(hashFunc, key)
	if _, err := hmacFunc.Write(input); err != nil {
		return nil, err
	}
	return hmacFunc.Sum(nil), nil
}

// Sha1ToHex Sign signs provided payload and returns encoded string sum.
func Sha1ToHex(data string) (string, error) {
	h := sha1.New() // #nosec
	if _, err := h.Write([]byte(data)); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

// HexEncodeToString takes in a hexadecimal byte array and returns a string
func HexEncodeToString(input []byte) string {
	return hex.EncodeToString(input)
}

// HexDecodeToBytes takes in a hexadecimal string and returns a byte array
func HexDecodeToBytes(input string) ([]byte, error) {
	return hex.DecodeString(input)
}

// ByteArrayToString returns a string
func ByteArrayToString(input []byte) string {
	return fmt.Sprintf("%x", input)
}

// Base64Decode takes in a Base64 string and returns a byte array and an error
func Base64Decode(input string) ([]byte, error) {
	result, err := base64.StdEncoding.DecodeString(input)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Base64Encode takes in a byte array then returns an encoded base64 string
func Base64Encode(input []byte) string {
	return base64.StdEncoding.EncodeToString(input)
}
