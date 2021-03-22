package env

import (
	"log"
	"os"
	"testing"

	"github.com/breathbath/go_utils/v2/pkg/errs"
	"github.com/stretchr/testify/assert"
)

func init() {
	err := os.Setenv("SOME_ENV_STRING", "GoLangCode")
	errs.FailOnError(err)

	err = os.Setenv("SOME_ENV_INT", "123")
	errs.FailOnError(err)

	err = os.Setenv("SOME_ENV_FLOAT", "123.345")
	errs.FailOnError(err)

	err = os.Setenv("SOME_ENV_BOOL_TRUE", "true")
	errs.FailOnError(err)

	err = os.Setenv("SOME_ENV_BOOL_FALSE", "false")
	errs.FailOnError(err)

	err = os.Setenv("SOME_ENV_BOOL_EMPTY", "")
	errs.FailOnError(err)
}

func TestSuite(t *testing.T) {
	defer func() {
		errCont := errs.NewErrorContainer()
		err := os.Unsetenv("SOME_ENV_STRING")
		errCont.AddError(err)

		err = os.Unsetenv("SOME_ENV_INT")
		errCont.AddError(err)

		err = os.Unsetenv("SOME_ENV_FLOAT")
		errCont.AddError(err)

		err = os.Unsetenv("SOME_ENV_BOOL_TRUE")
		errCont.AddError(err)

		err = os.Unsetenv("SOME_ENV_BOOL_FALSE")
		errCont.AddError(err)

		err = os.Unsetenv("SOME_ENV_BOOL_EMPTY")
		errCont.AddError(err)

		err = errCont.Result("\n")
		if err != nil {
			log.Println(err)
		}
	}()
	t.Run("testReadEnv", testReadEnv)
	t.Run("testReadEnvInt64", testReadEnvInt64)
	t.Run("testReadEnvInt", testReadEnvInt)
	t.Run("testReadEnvFloat", testReadEnvFloat)
	t.Run("testReadEnvOrError", testReadEnvOrError)
	t.Run("testReadEnvOrFailWithExistingEnv", testReadEnvOrFailWithExistingEnv)
	t.Run("testReadEnvOrFailWithNonExistingEnv", testReadEnvOrFailWithNonExistingEnv)
	t.Run("testReadEnvBool", testReadEnvBool)
}

func testReadEnv(t *testing.T) {
	assert.Equal(t, "GoLangCode", ReadEnv("SOME_ENV_STRING", ""))
	assert.Equal(t, "Some default val", ReadEnv("NON_EXISTING_ENV_STRING", "Some default val"))
}

func testReadEnvInt64(t *testing.T) {
	assert.EqualValues(t, 123, ReadEnvInt64("SOME_ENV_INT", 0))
	assert.EqualValues(t, 0, ReadEnvInt64("NON_EXISTING_ENV_INT", 0))
	assert.EqualValues(t, 2, ReadEnvInt64("SOME_ENV_STRING", 2))
}

func testReadEnvBool(t *testing.T) {
	assert.EqualValues(t, true, ReadEnvBool("SOME_ENV_BOOL_TRUE", false))
	assert.EqualValues(t, false, ReadEnvBool("SOME_ENV_BOOL_FALSE", true))
	assert.EqualValues(t, false, ReadEnvBool("SOME_ENV_BOOL_EMPTY", true))
	assert.EqualValues(t, false, ReadEnvBool("NON_EXISTING_ENV_BOOL", false))
	assert.EqualValues(t, true, ReadEnvBool("NON_EXISTING_ENV_BOOL", true))
}

func testReadEnvInt(t *testing.T) {
	assert.EqualValues(t, 123, ReadEnvInt("SOME_ENV_INT", 0))
	assert.EqualValues(t, 0, ReadEnvInt("NON_EXISTING_ENV_INT", 0))
	assert.EqualValues(t, 0, ReadEnvInt("SOME_ENV_STRING", 0))
}

func testReadEnvFloat(t *testing.T) {
	assert.EqualValues(t, 123.345, ReadEnvFloat("SOME_ENV_FLOAT", 0.1))
	assert.EqualValues(t, 0.1, ReadEnvFloat("NON_EXISTING_ENV_INT", 0.1))
	assert.EqualValues(t, 0.2, ReadEnvFloat("SOME_ENV_STRING", 0.2))
}

func testReadEnvOrError(t *testing.T) {
	actualVal, err := ReadEnvOrError("SOME_ENV_STRING")
	assert.NoError(t, err)
	assert.Equal(t, "GoLangCode", actualVal)

	actualVal, err = ReadEnvOrError("NON_EXISTING_ENV_STRING")
	assert.EqualError(t, err, "required env variable 'NON_EXISTING_ENV_STRING' is not set")
	assert.Equal(t, "", actualVal)
}

func testReadEnvOrFailWithExistingEnv(t *testing.T) {
	actualVal := ReadEnvOrFail("SOME_ENV_STRING")
	assert.Equal(t, "GoLangCode", actualVal)
}

func testReadEnvOrFailWithNonExistingEnv(t *testing.T) {
	assert.PanicsWithValue(t, "required env variable 'NON_EXISTING_ENV_STRING' is not set", func() {
		ReadEnvOrFail("NON_EXISTING_ENV_STRING")
	})
}
