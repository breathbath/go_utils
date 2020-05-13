package options

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"os"
	"strings"
	"testing"
)

func TestMapValuesProvider(t *testing.T) {
	mvp := NewMapValuesProvider(map[string]interface{}{
		"one": 1,
		"two": "2",
	})

	val, found := mvp.Read("one")
	assert.True(t, found)
	assert.Equal(t, 1, val)

	val2, found2 := mvp.Read("two")
	assert.True(t, found2)
	assert.Equal(t, "2", val2)

	val3, found3 := mvp.Read("three")
	assert.False(t, found3)
	assert.Equal(t, nil, val3)

	mvp2 := mvp.Copy(map[string]interface{}{
		"four": true,
	})
	val4, found4 := mvp2.Read("one")
	assert.True(t, found4)
	assert.Equal(t, 1, val4)

	val5, found5 := mvp2.Read("four")
	assert.True(t, found5)
	assert.Equal(t, true, val5)
}

func TestEnvValuesProvider(t *testing.T) {
	err := os.Setenv("someenv", "someenvval")
	assert.NoError(t, err)

	evp := EnvValuesProvider{}

	val, found := evp.Read("someenv")
	assert.True(t, found)
	assert.Equal(t, "someenvval", val)

	_, found2 := evp.Read("someindff")
	assert.False(t, found2)
}

func TestJsonFileValuesProvider(t *testing.T) {
	buf := strings.NewReader(`{"key1":"val1","key2":2,"key3":3.3,"key4":null,"key5":""}`)
	jvp, err := NewJsonValuesProvider(buf)
	assert.NoError(t, err)

	val, found := jvp.Read("key1")
	assert.True(t, found)
	assert.Equal(t, "val1", val)

	val2, found2 := jvp.Read("key2")
	assert.True(t, found2)
	assert.Equal(t, 2, val2)

	val3, found3 := jvp.Read("key3")
	assert.True(t, found3)
	assert.Equal(t, 3.3, val3)

	val4, found4 := jvp.Read("key4")
	assert.True(t, found4)
	assert.Equal(t, nil, val4)

	val5, found5 := jvp.Read("key5")
	assert.True(t, found5)
	assert.Equal(t, "", val5)

	_, found6 := jvp.Read("someNonExistingVal")
	assert.False(t, found6)
}

type failingReader struct{}
func (fr failingReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("some error")
}

func TestJsonFileValuesProviderFailedReader(t *testing.T) {
	_, err := NewJsonValuesProvider(failingReader{})
	assert.EqualError(t, err, "some error")
}

func TestJsonFileValuesProviderInvalidJson(t *testing.T) {
	buf := strings.NewReader(`dfasdf`)
	_, err := NewJsonValuesProvider(buf)
	assert.Error(t, err)
}

func TestValuesProviderComposite(t *testing.T) {
	mvp1 := NewMapValuesProvider(map[string]interface{}{
		"one": 1,
		"two": "2",
	})
	mvp2 := NewMapValuesProvider(map[string]interface{}{
		"three": 3,
	})
	vpc := NewValuesProviderComposite(mvp1, mvp2)

	val, found := vpc.Read("one")
	assert.True(t, found)
	assert.Equal(t, 1, val)

	val2, found2 := vpc.Read("two")
	assert.True(t, found2)
	assert.Equal(t, "2", val2)

	val4, found4 := vpc.Read("three")
	assert.True(t, found4)
	assert.Equal(t, 3, val4)

	_, found3 := vpc.Read("notFound")
	assert.False(t, found3)
}

func TestParameterBag(t *testing.T) {
	mvp := NewMapValuesProvider(map[string]interface{}{
		"intval": 1,
		"strval": "someStr",
		"struct": struct{a bool}{a: true},
	})

	pb := New(mvp)
	val, found := pb.Read("intval", 12)
	assert.True(t, found)
	assert.Equal(t, 1, val)

	val2, found2 := pb.Read("notfoundval", 12)
	assert.False(t, found2)
	assert.Equal(t, 12, val2)

	val3, err := pb.ReadRequired("intval")
	assert.NoError(t, err)
	assert.Equal(t, 1, val3)

	_, err = pb.ReadRequired("notfoundval")
	assert.EqualError(t, err, "required option notfoundval is empty")

	val4 := pb.ReadString("intval", "12")
	assert.Equal(t, "1", val4)

	val5 := pb.ReadString("notfoundval", "defaultVal")
	assert.Equal(t, "defaultVal", val5)

	val55 := pb.ReadString("struct", "someVal")
	assert.Equal(t, "{true}", val55)

	val6, err := pb.ReadRequiredString("strval")
	assert.NoError(t, err)
	assert.Equal(t, "someStr", val6)

	_, err = pb.ReadRequiredString("notfoundval")
	assert.EqualError(t, err, "required option notfoundval is empty")

	val7, err := pb.ReadRequiredString("struct")
	assert.NoError(t, err,)
	assert.Equal(t, "{true}", val7)
}
