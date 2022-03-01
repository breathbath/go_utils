package enc

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"hash/crc32"
	"strconv"
)

func Checksum(s string) string {
	crc32InUint32 := crc32.ChecksumIEEE([]byte(s))
	crc32InString := strconv.FormatUint(uint64(crc32InUint32), 16)
	return crc32InString
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
