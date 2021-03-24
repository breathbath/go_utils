package options

import (
	"bytes"
	"errors"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMapValuesProviderRead(t *testing.T) {
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

func TestMapValuesProviderDump(t *testing.T) {
	mvp := NewMapValuesProvider(map[string]interface{}{
		"one": 1,
		"two": "2",
	})

	b := &bytes.Buffer{}
	err := mvp.Dump(b)
	assert.NoError(t, err)
	if err != nil {
		return
	}

	assert.Contains(t, b.String(), `"one":1`)
	assert.Contains(t, b.String(), `"two":"2"`)
}

func TestMapValuesProviderToKeyValue(t *testing.T) {
	providedKeyValues := map[string]interface{}{
		"one": 1,
		"two": "2",
	}
	mvp := NewMapValuesProvider(providedKeyValues)

	actualKeyValues := mvp.ToKeyValues()

	assert.Equal(t, providedKeyValues, actualKeyValues)
}

func TestEnvValuesProviderRead(t *testing.T) {
	err := os.Setenv("someenv", "someenvval")
	assert.NoError(t, err)
	if err != nil {
		return
	}
	defer func() {
		err = os.Unsetenv("someenv")
		if err != nil {
			log.Println(err.Error())
		}
	}()

	evp := EnvValuesProvider{}

	val, found := evp.Read("someenv")
	assert.True(t, found)
	assert.Equal(t, "someenvval", val)

	_, found2 := evp.Read("someindff")
	assert.False(t, found2)
}

func TestEnvValuesProviderDump(t *testing.T) {
	err := os.Setenv("someenv1", "someenvval1")
	assert.NoError(t, err)
	if err != nil {
		return
	}
	defer func() {
		err = os.Unsetenv("someenv1")
		if err != nil {
			log.Println(err.Error())
		}
	}()

	evp := EnvValuesProvider{}

	b := &bytes.Buffer{}
	err = evp.Dump(b)
	assert.NoError(t, err)
	if err != nil {
		return
	}

	assert.Contains(t, b.String(), `someenv1`)
	assert.Contains(t, b.String(), `someenvval1`)
}

func TestEnvValuesProviderToKeyValues(t *testing.T) {
	err := os.Setenv("someenv222", "someenvval222")
	assert.NoError(t, err)
	if err != nil {
		return
	}
	defer func() {
		err = os.Unsetenv("someenv222")
		if err != nil {
			log.Println(err.Error())
		}
	}()

	evp := EnvValuesProvider{}

	kvs := evp.ToKeyValues()

	assert.Equal(t, "someenvval222", kvs["someenv222"])
}

func TestJsonFileValuesProviderRead(t *testing.T) {
	buf := strings.NewReader(`{"key1":"val1","key2":2,"key3":3.3,"key4":null,"key5":""}`)
	jvp, err := NewJSONValuesProvider(buf)
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

func TestJsonFileValuesProviderDump(t *testing.T) {
	jsonStr := `{"key1":"val1","key2":2,"key3":3.3,"key4":null,"key5":""}`
	readBuf := strings.NewReader(jsonStr)
	jvp, err := NewJSONValuesProvider(readBuf)
	assert.NoError(t, err)
	if err != nil {
		return
	}

	writeBuf := &bytes.Buffer{}
	err = jvp.Dump(writeBuf)
	assert.NoError(t, err)
	if err != nil {
		return
	}

	actualDumpValue := writeBuf.String()
	assert.Contains(t, actualDumpValue, `"key1":"val1"`)
	assert.Contains(t, actualDumpValue, `"key2":2,"key3"`)
	assert.Contains(t, actualDumpValue, `"key4":null`)
	assert.Contains(t, actualDumpValue, `"key5":""`)
}

func TestJsonFileValuesProviderToKeyValues(t *testing.T) {
	jsonStr := `{"color":"red","number":2}`
	readBuf := strings.NewReader(jsonStr)
	jvp, err := NewJSONValuesProvider(readBuf)
	assert.NoError(t, err)
	if err != nil {
		return
	}

	kvs := jvp.ToKeyValues()

	assert.Equal(t, map[string]interface{}{"color": "red", "number": 2}, kvs)
}

type failingReader struct{}

func (fr failingReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("some error")
}

func TestJsonFileValuesProviderFailedReader(t *testing.T) {
	_, err := NewJSONValuesProvider(failingReader{})
	assert.EqualError(t, err, "some error")
}

func TestJsonFileValuesProviderInvalidJson(t *testing.T) {
	buf := strings.NewReader(`dfasdf`)
	_, err := NewJSONValuesProvider(buf)
	assert.Error(t, err)
}

func TestValuesProviderCompositeRead(t *testing.T) {
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

func TestValuesProviderCompositeDump(t *testing.T) {
	mvp1 := NewMapValuesProvider(map[string]interface{}{
		"name":    "John",
		"surname": "Deer",
	})
	mvp2 := NewMapValuesProvider(map[string]interface{}{
		"age": 33,
	})
	vpc := NewValuesProviderComposite(mvp1, mvp2)

	writeBuf := &bytes.Buffer{}
	err := vpc.Dump(writeBuf)
	assert.NoError(t, err)
	if err != nil {
		return
	}

	actualDumpValue := writeBuf.String()
	assert.Contains(t, actualDumpValue, `"name":"John"`)
	assert.Contains(t, actualDumpValue, `"surname":"Deer"`)
	assert.Contains(t, actualDumpValue, `"age":33`)
}

func TestValuesProviderCompositeToKeyValues(t *testing.T) {
	jsonStr := `{"make":"four","take":"three"}`
	readBuf := strings.NewReader(jsonStr)
	jvp, err := NewJSONValuesProvider(readBuf)
	assert.NoError(t, err)
	if err != nil {
		return
	}

	err = os.Setenv("onekey", "oneValue")
	assert.NoError(t, err)
	if err != nil {
		return
	}
	defer func() {
		err = os.Unsetenv("onekey")
		if err != nil {
			log.Println(err.Error())
		}
	}()

	evp := &EnvValuesProvider{}

	vpc := NewValuesProviderComposite(jvp, evp)
	kvs := vpc.ToKeyValues()

	assert.Equal(t, "four", kvs["make"])
	assert.Equal(t, "three", kvs["take"])
	assert.Equal(t, "oneValue", kvs["onekey"])
}

func TestParameterBagRead(t *testing.T) {
	mvp := NewMapValuesProvider(map[string]interface{}{
		"intval":          1,
		"intval64":        int64(12),
		"strval":          "someStr",
		"strvalEmpty":     "",
		"struct":          struct{ a bool }{a: true},
		"stringsVal":      []string{"one", "two"},
		"stringsValEmpty": []string{},
		"intStrVal":       "33",
		"dur":             time.Second * 34,
		"boolStrFalse":    "false",
		"boolStrTrue":     "true",
		"boolFalse":       false,
		"boolTrue":        true,
		"int0":            0,
		"int-1":           -1,
		"str0":            "0",
		"unitVal":         uint(14),
		"nil":             nil,
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

	val551 := pb.ReadString("strval", "someVal")
	assert.Equal(t, "someStr", val551)

	val6, err := pb.ReadRequiredString("strval")
	assert.NoError(t, err)
	assert.Equal(t, "someStr", val6)

	_, err = pb.ReadRequiredString("notfoundval")
	assert.EqualError(t, err, "required option notfoundval is empty")

	val7, err := pb.ReadRequiredString("struct")
	assert.NoError(t, err)
	assert.Equal(t, "{true}", val7)

	_, err = pb.ReadRequiredString("strvalEmpty")
	assert.EqualError(t, err, "required option strvalEmpty is empty")

	val8 := pb.ReadStrings("stringsVal", "three")
	assert.Equal(t, []string{"one", "two"}, val8)

	val9 := pb.ReadStrings("notfoundval", "three", "four")
	assert.Equal(t, []string{"three", "four"}, val9)

	val10 := pb.ReadStrings("strval", "three")
	assert.Equal(t, []string{"someStr"}, val10)

	val11 := pb.ReadStrings("struct", "three")
	assert.Equal(t, []string{"three"}, val11)

	val12, err := pb.ReadRequiredStrings("stringsVal")
	assert.NoError(t, err)
	assert.Equal(t, []string{"one", "two"}, val12)

	val13, err := pb.ReadRequiredStrings("strval")
	assert.NoError(t, err)
	assert.Equal(t, []string{"someStr"}, val13)

	_, err = pb.ReadRequiredStrings("notfoundval")
	assert.EqualError(t, err, "required option notfoundval is empty")

	_, err = pb.ReadRequiredStrings("struct")
	assert.EqualError(t, err, "cannot convert value {true} to []string")

	val14 := pb.ReadInt("intval", 14)
	assert.Equal(t, 1, val14)

	val15 := pb.ReadInt("notfoundval", 15)
	assert.Equal(t, 15, val15)

	val555 := pb.ReadInt("struct", 14)
	assert.Equal(t, 14, val555)

	val556 := pb.ReadInt("intStrVal", 22)
	assert.Equal(t, 33, val556)

	val16, err := pb.ReadRequiredInt("intval")
	assert.NoError(t, err)
	assert.Equal(t, 1, val16)

	_, err = pb.ReadRequiredInt("notfoundval")
	assert.EqualError(t, err, "required option notfoundval is empty")

	val17, err := pb.ReadRequiredInt("intStrVal")
	assert.NoError(t, err)
	assert.Equal(t, 33, val17)

	_, err = pb.ReadRequiredInt("strval")
	assert.EqualError(t, err, "cannot convert someStr to int")

	_, err = pb.ReadRequiredInt("struct")
	assert.EqualError(t, err, "cannot convert {true} to int")

	val18 := pb.ReadInt64("intval", 14)
	assert.Equal(t, int64(1), val18)

	val19 := pb.ReadInt64("notfoundval", 15)
	assert.Equal(t, int64(15), val19)

	val20 := pb.ReadInt64("struct", 14)
	assert.Equal(t, int64(14), val20)

	val21 := pb.ReadInt64("intStrVal", 22)
	assert.Equal(t, int64(33), val21)

	val211 := pb.ReadInt64("intval64", 22)
	assert.Equal(t, int64(12), val211)

	val22, err := pb.ReadRequiredInt64("intval")
	assert.NoError(t, err)
	assert.Equal(t, int64(1), val22)

	val221, err := pb.ReadRequiredInt64("intval64")
	assert.NoError(t, err)
	assert.Equal(t, int64(12), val221)

	_, err = pb.ReadRequiredInt64("notfoundval")
	assert.EqualError(t, err, "required option notfoundval is empty")

	val23, err := pb.ReadRequiredInt64("intStrVal")
	assert.NoError(t, err)
	assert.Equal(t, int64(33), val23)

	_, err = pb.ReadRequiredInt64("strval")
	assert.EqualError(t, err, "cannot convert someStr to int64")

	_, err = pb.ReadRequiredInt64("struct")
	assert.EqualError(t, err, "cannot convert {true} to int64")

	val24 := pb.ReadDuration("intval", time.Second, 55)
	assert.Equal(t, time.Second, val24)

	val25 := pb.ReadDuration("notfoundval", time.Second, 55)
	assert.Equal(t, time.Second*55, val25)

	val26 := pb.ReadDuration("struct", time.Second, 55)
	assert.Equal(t, time.Second*55, val26)

	val27 := pb.ReadDuration("intStrVal", time.Second, 56)
	assert.Equal(t, time.Second*33, val27)

	val28 := pb.ReadDuration("dur", time.Second, 57)
	assert.Equal(t, time.Second*34, val28)

	val29, err := pb.ReadRequiredDuration("dur", time.Minute)
	assert.NoError(t, err)
	assert.Equal(t, time.Second*34, val29)

	val30, err := pb.ReadRequiredDuration("intStrVal", time.Minute)
	assert.NoError(t, err)
	assert.Equal(t, time.Minute*33, val30)

	_, err = pb.ReadRequiredDuration("notfoundval", time.Minute)
	assert.EqualError(t, err, "required option notfoundval is empty")

	val31, err := pb.ReadRequiredDuration("intval", time.Minute)
	assert.NoError(t, err)
	assert.Equal(t, time.Minute, val31)

	_, err = pb.ReadRequiredDuration("strval", time.Minute)
	assert.EqualError(t, err, "cannot convert someStr to uint")

	_, err = pb.ReadRequiredDuration("struct", time.Minute)
	assert.EqualError(t, err, "cannot convert {true} to uint")

	val32 := pb.ReadBool("boolStrFalse", true)
	assert.Equal(t, false, val32)

	val33 := pb.ReadBool("notfoundval", true)
	assert.Equal(t, true, val33)

	val34 := pb.ReadBool("notfoundval", false)
	assert.Equal(t, false, val34)

	val35 := pb.ReadBool("int0", true)
	assert.Equal(t, false, val35)

	val36 := pb.ReadBool("str0", true)
	assert.Equal(t, false, val36)

	val37 := pb.ReadBool("strvalEmpty", true)
	assert.Equal(t, false, val37)

	val38 := pb.ReadBool("stringsValEmpty", true)
	assert.Equal(t, false, val38)

	val39 := pb.ReadBool("intval", false)
	assert.Equal(t, true, val39)

	val40 := pb.ReadBool("intval64", false)
	assert.Equal(t, true, val40)

	val41 := pb.ReadBool("strval", false)
	assert.Equal(t, true, val41)

	val42 := pb.ReadBool("stringsVal", false)
	assert.Equal(t, true, val42)

	val43 := pb.ReadBool("dur", false)
	assert.Equal(t, true, val43)

	val44 := pb.ReadBool("boolStrTrue", false)
	assert.Equal(t, true, val44)

	val441 := pb.ReadBool("boolTrue", false)
	assert.Equal(t, true, val441)

	val442 := pb.ReadBool("boolFalse", true)
	assert.Equal(t, false, val442)

	val45, err := pb.ReadRequiredBool("boolStrTrue")
	assert.NoError(t, err)
	assert.Equal(t, true, val45)

	val451, err := pb.ReadRequiredBool("boolTrue")
	assert.NoError(t, err)
	assert.Equal(t, true, val451)

	val452, err := pb.ReadRequiredBool("boolFalse")
	assert.NoError(t, err)
	assert.Equal(t, false, val452)

	val46, err := pb.ReadRequiredBool("int0")
	assert.NoError(t, err)
	assert.Equal(t, false, val46)

	_, err = pb.ReadRequiredBool("notfoundval")
	assert.EqualError(t, err, "required option notfoundval is empty")

	val47, err := pb.ReadRequiredBool("str0")
	assert.NoError(t, err)
	assert.Equal(t, false, val47)

	val48 := pb.ReadUint("intval", 14)
	assert.Equal(t, uint(1), val48)

	val49 := pb.ReadUint("unitVal", 14)
	assert.Equal(t, uint(14), val49)

	val50 := pb.ReadUint("notfoundval", 15)
	assert.Equal(t, uint(15), val50)

	val51 := pb.ReadUint("struct", 14)
	assert.Equal(t, uint(14), val51)

	val52 := pb.ReadUint("intStrVal", 22)
	assert.Equal(t, uint(33), val52)

	val53 := pb.ReadUint("intval64", 22)
	assert.Equal(t, uint(12), val53)

	val54 := pb.ReadUint("int-1", 23)
	assert.Equal(t, uint(23), val54)

	val56, err := pb.ReadRequiredUint("unitVal")
	assert.NoError(t, err)
	assert.Equal(t, uint(14), val56)

	_, err = pb.ReadRequiredUint("int-1")
	assert.EqualError(t, err, "cannot convert -1 to uint")

	_, err = pb.ReadRequiredUint("notfoundval")
	assert.EqualError(t, err, "required option notfoundval is empty")

	val57, err := pb.ReadRequiredUint("intStrVal")
	assert.NoError(t, err)
	assert.Equal(t, uint(33), val57)

	_, err = pb.ReadRequiredUint("strval")
	assert.EqualError(t, err, "cannot convert someStr to uint")

	_, err = pb.ReadRequiredUint("struct")
	assert.EqualError(t, err, "cannot convert {true} to uint")

	err = pb.CheckRequiredValues([]string{"notfoundVal", "boolFalse", "notFoundVAl3"})
	assert.EqualError(t, err, "required option notfoundVal is empty required option notFoundVAl3 is empty")

	err = pb.CheckRequiredValues([]string{"intval", "boolFalse", "intval64", "stringsValEmpty", "strvalEmpty", "str0", "nil"})
	assert.NoError(t, err)
}

func TestParameterBagWithNoValuesProvider(t *testing.T) {
	pb := New(nil)
	val := pb.ReadString("someKey", "someDefault")
	assert.Equal(t, "someDefault", val)
}

func TestParameterBagMerge(t *testing.T) {
	mvp1 := NewMapValuesProvider(map[string]interface{}{
		"randomNumb": 134,
	})
	mvp2 := NewMapValuesProvider(map[string]interface{}{
		"randomStr": "lsls4",
	})

	paramBag1 := New(mvp1)
	paramBag2 := New(mvp2)

	paramBag1.MergeParameterBag(paramBag2)

	assert.Equal(t, 134, paramBag1.ReadInt("randomNumb", 0))
	assert.Equal(t, "lsls4", paramBag1.ReadString("randomStr", ""))
}
