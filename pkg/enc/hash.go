package enc

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"hash/crc32"
	"strconv"

	uuid "github.com/nu7hatch/gouuid"
	"golang.org/x/crypto/bcrypt"
)

func Hash(s string) string {
	hasher := md5.New()
	_, err := hasher.Write([]byte(s))
	if err != nil {
		return ""
	}

	return hex.EncodeToString(hasher.Sum(nil))
}

func Checksum(s string) string {
	crc32InUint32 := crc32.ChecksumIEEE([]byte(s))
	crc32InString := strconv.FormatUint(uint64(crc32InUint32), 16)
	return crc32InString
}

func NewUUID(s string) string {
	h := Hash(s)
	u5, _ := uuid.NewV5(uuid.NamespaceURL, []byte(h))
	return u5.String()
}

func GetPasswordHash(input string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input), 10)

	return string(hashedPassword), err
}

func CreateMd5Hash(key string) string {
	hasher := md5.New()
	_, err := hasher.Write([]byte(key))
	if err != nil {
		return ""
	}
	return hex.EncodeToString(hasher.Sum(nil))
}

func ComputeHmac256Hex(message, secret string) string {
	key := []byte(secret)
	sig := hmac.New(sha256.New, key)

	_, err := sig.Write([]byte(message))
	if err != nil {
		return ""
	}

	return hex.EncodeToString(sig.Sum(nil))
}

func ComputeHmac256Base64(message string, secret []byte) string {
	sig := hmac.New(sha256.New, secret)
	_, err := sig.Write([]byte(message))
	if err != nil {
		return ""
	}

	base64EncodedSignature := base64.StdEncoding.EncodeToString(sig.Sum(nil))

	return base64EncodedSignature
}
