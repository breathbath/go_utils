package enc

import (
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func TestHash(t *testing.T) {
	assert.Equal(t, "bcbe3365e6ac95ea2c0343a2395834dd", Hash("222"))
}

func TestChecksum(t *testing.T) {
	assert.Equal(t, "fd09ed1d", Checksum("222"))
}

func TestNewUuid(t *testing.T) {
	assert.Equal(t, "d28b7a82-fe8f-57af-43e1-17bada1fd734", NewUuid("222"))
}

func TestGetPasswordHash(t *testing.T) {
	actualOutput, err := GetPasswordHash("1234567890")
	assert.NoError(t, err)
	err = bcrypt.CompareHashAndPassword([]byte(actualOutput), []byte("1234567890"))
	assert.NoError(t, err)
}

func TestCreateMd5Hash(t *testing.T) {
	assert.Equal(t, "bcbe3365e6ac95ea2c0343a2395834dd", CreateMd5Hash("222"))
}

func TestComputeHmac256Hex(t *testing.T) {
	assert.Equal(t, "446dc894180cbae72aff1988d2c5a595e23ca05f331b5a411874e1da7a159044", ComputeHmac256Hex("some msg", "123"))
}

func TestComputeHmac256Base64(t *testing.T) {
	assert.Equal(t, "RG3IlBgMuucq/xmI0sWlleI8oF8zG1pBGHTh2noVkEQ=", ComputeHmac256Base64("some msg", []byte("123")))
}
