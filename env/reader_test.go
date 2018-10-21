package env

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func init() {
	os.Setenv("SOME_ENV_STRING", "GoLangCode")
	os.Setenv("SOME_ENV_INT", "123")
	os.Setenv("SOME_ENV_FLOAT", "123.345")
}

func TestReadEnv(t *testing.T) {
	assert.Equal(t, "GoLangCode", ReadEnv("SOME_ENV_STRING", ""))
	assert.Equal(t, "Some default val", ReadEnv("NON_EXISTING_ENV_STRING", "Some default val"))
}

func TestReadEnvInt(t *testing.T) {
	assert.EqualValues(t, 123, ReadEnvInt("SOME_ENV_INT", 0))
	assert.EqualValues(t, 0, ReadEnvInt("NON_EXISTING_ENV_INT", 0))
}

func TestReadEnvFloat(t *testing.T) {
	assert.EqualValues(t, 123.345, ReadEnvFloat("SOME_ENV_FLOAT", 0.1))
	assert.EqualValues(t, 0.1, ReadEnvFloat("NON_EXISTING_ENV_INT", 0.1))
}

func TestReadEnvOrError(t *testing.T) {
	actualVal, err := ReadEnvOrError("SOME_ENV_STRING")
	assert.NoError(t, err)
	assert.Equal(t, "GoLangCode", actualVal)

	actualVal, err = ReadEnvOrError("NON_EXISTING_ENV_STRING")
	assert.EqualError(t, err, "Required env variable 'NON_EXISTING_ENV_STRING' is not set")
	assert.Equal(t, "", actualVal)
}

func TestReadEnvOrFailWithExistingEnv(t *testing.T) {
	actualVal := ReadEnvOrFail("SOME_ENV_STRING")
	assert.Equal(t, "GoLangCode", actualVal)
}

func TestReadEnvOrFailWithNonExistingEnv(t *testing.T) {
	assert.PanicsWithValue(t, "Required env variable 'NON_EXISTING_ENV_STRING' is not set", func() {
		ReadEnvOrFail("NON_EXISTING_ENV_STRING")
	})
}