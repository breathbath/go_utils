package enc

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"github.com/nu7hatch/gouuid"
	"golang.org/x/crypto/bcrypt"
	"hash/crc32"
	"strconv"
)

func Hash(s string) string {
	hasher := md5.New()
	hasher.Write([]byte(s))
	return hex.EncodeToString(hasher.Sum(nil))
}

func Checksum(s string) string {
	crc32InUint32 := crc32.ChecksumIEEE([]byte(s))
	crc32InString := strconv.FormatUint(uint64(crc32InUint32), 16)
	return crc32InString
}

func NewUuid(s string) string {
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
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}

func ComputeHmac256Hex(message string, secret string) string {
	key := []byte(secret)
	sig := hmac.New(sha256.New, key)
	sig.Write([]byte(message))

	return hex.EncodeToString(sig.Sum(nil))
}

func ComputeHmac256Base64(message string, secret []byte) string {
	sig := hmac.New(sha256.New, secret)
	sig.Write([]byte(message))
	base64EncodedSignature := base64.StdEncoding.EncodeToString(sig.Sum(nil))

	return base64EncodedSignature
}
