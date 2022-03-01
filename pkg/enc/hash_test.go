package enc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChecksum(t *testing.T) {
	assert.Equal(t, "fd09ed1d", Checksum("222"))
}

func TestComputeHmac256Hex(t *testing.T) {
	assert.Equal(t, "446dc894180cbae72aff1988d2c5a595e23ca05f331b5a411874e1da7a159044", ComputeHmac256Hex("some msg", "123"))
}

func TestComputeHmac256Base64(t *testing.T) {
	assert.Equal(t, "RG3IlBgMuucq/xmI0sWlleI8oF8zG1pBGHTh2noVkEQ=", ComputeHmac256Base64("some msg", []byte("123")))
}
