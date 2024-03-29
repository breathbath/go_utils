package options

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/breathbath/go_utils/v3/pkg/conv"
	errs2 "github.com/breathbath/go_utils/v3/pkg/errs"
)

type ValuesProvider interface {
	Read(name string) (val interface{}, found bool)
	Dump(w io.Writer) (err error)
	ToKeyValues() map[string]interface{}
}

type MapValuesProvider struct {
	parameters *sync.Map
}

func NewMapValuesProvider(params map[string]interface{}) *MapValuesProvider {
	paramsMap := &sync.Map{}
	for key, val := range params {
		paramsMap.Store(key, val)
	}

	bag := &MapValuesProvider{
		parameters: paramsMap,
	}

	return bag
}

// Copy takes internal items, merges them with newParams and returns the result
func (mvp *MapValuesProvider) Copy(newParams map[string]interface{}) *MapValuesProvider {
	resultItems := &sync.Map{}
	mvp.parameters.Range(func(key, value interface{}) bool {
		resultItems.Store(key, value)
		return true
	})

	for k, v := range newParams {
		resultItems.Store(k, v)
	}

	return &MapValuesProvider{
		parameters: resultItems,
	}
}

func (mvp *MapValuesProvider) Read(name string) (val interface{}, found bool) {
	return mvp.parameters.Load(name)
}

func (mvp *MapValuesProvider) ToKeyValues() map[string]interface{} {
	return conv.ConvertSyncMapToMap(mvp.parameters)
}

func (mvp *MapValuesProvider) Dump(w io.Writer) (err error) {
	data := conv.ConvertSyncMapToMap(mvp.parameters)
	jsonEncoder := json.NewEncoder(w)
	err = jsonEncoder.Encode(data)
	return
}

type NullValuesProvider struct{}

func (nvp *NullValuesProvider) Read(name string) (val interface{}, found bool) {
	return
}

func (nvp *NullValuesProvider) Dump(w io.Writer) (err error) {
	return
}

func (nvp *NullValuesProvider) ToKeyValues() map[string]interface{} {
	return map[string]interface{}{}
}

type EnvValuesProvider struct{}

func (evp *EnvValuesProvider) Read(name string) (val interface{}, found bool) {
	return os.LookupEnv(name)
}

func (evp *EnvValuesProvider) Dump(w io.Writer) (err error) {
	data := os.Environ()
	jsonEncoder := json.NewEncoder(w)
	err = jsonEncoder.Encode(data)
	return
}

func (evp *EnvValuesProvider) ToKeyValues() map[string]interface{} {
	res := map[string]interface{}{}
	for _, env := range os.Environ() {
		envPair := strings.SplitN(env, "=", 2)
		res[envPair[0]] = envPair[1]
	}

	return res
}

type JSONFileValuesProvider struct {
	vals MapValuesProvider
}

func NewJSONValuesProvider(jsond io.Reader) (jfvp *JSONFileValuesProvider, err error) {
	var data []byte
	data, err = io.ReadAll(jsond)
	if err != nil {
		return
	}

	var objmap map[string]json.RawMessage
	err = json.Unmarshal(data, &objmap)
	if err != nil {
		return
	}

	params := &sync.Map{}
	for k, val := range objmap {
		valStr := string(val)
		if strings.HasPrefix(valStr, `"`) || valStr == "" {
			params.Store(k, strings.Trim(string(val), `"`))
			continue
		}

		if valStr == "null" {
			params.Store(k, nil)
			continue
		}

		intVal, err := strconv.Atoi(valStr)
		if err == nil {
			params.Store(k, intVal)
			continue
		}

		floatVal, err := strconv.ParseFloat(valStr, 64)
		if err == nil {
			params.Store(k, floatVal)
			continue
		}

		params.Store(k, valStr)
	}

	vals := MapValuesProvider{
		parameters: params,
	}

	return &JSONFileValuesProvider{vals: vals}, nil
}

func (jfvp *JSONFileValuesProvider) Read(name string) (val interface{}, found bool) {
	return jfvp.vals.Read(name)
}

func (jfvp *JSONFileValuesProvider) Dump(w io.Writer) (err error) {
	return jfvp.vals.Dump(w)
}

func (jfvp *JSONFileValuesProvider) ToKeyValues() map[string]interface{} {
	return jfvp.vals.ToKeyValues()
}

type ValuesProviderComposite struct {
	providers []ValuesProvider
}

func NewValuesProviderComposite(vps ...ValuesProvider) *ValuesProviderComposite {
	return &ValuesProviderComposite{providers: vps}
}

func (vpc *ValuesProviderComposite) Read(name string) (val interface{}, found bool) {
	for _, vp := range vpc.providers {
		val, found = vp.Read(name)
		if found {
			return
		}
	}

	return
}

func (vpc *ValuesProviderComposite) ToKeyValues() map[string]interface{} {
	res := map[string]interface{}{}
	for _, vp := range vpc.providers {
		kvs := vp.ToKeyValues()
		for k, v := range kvs {
			res[k] = v
		}
	}

	return res
}

func (vpc *ValuesProviderComposite) Dump(w io.Writer) (err error) {
	kvs := vpc.ToKeyValues()
	jsonEncoder := json.NewEncoder(w)
	return jsonEncoder.Encode(kvs)
}

// ParameterBag construction for holding configuration options
type ParameterBag struct {
	BaseValuesProvider ValuesProvider
}

// New creates empty bag
func New(vp ValuesProvider) *ParameterBag {
	if vp == nil {
		vp = &NullValuesProvider{}
	}

	return &ParameterBag{
		BaseValuesProvider: vp,
	}
}

// Read reads interface value, if not found, will read from envs, if not found there will return defaultVal
func (p *ParameterBag) Read(name string, defaultVal interface{}) (interface{}, bool) {
	if p.BaseValuesProvider == nil {
		p.BaseValuesProvider = &NullValuesProvider{}
	}
	val, found := p.BaseValuesProvider.Read(name)
	if found {
		return val, found
	}

	return defaultVal, found
}

// ReadRequired reads interface value, if not found, will return error
func (p *ParameterBag) ReadRequired(name string) (interface{}, error) {
	valI, found := p.Read(name, nil)
	if !found {
		return valI, fmt.Errorf("required option %s is empty", name)
	}

	return valI, nil
}

// ReadString same as Read but returns a string
func (p *ParameterBag) ReadString(name, defaultVal string) string {
	valI, found := p.Read(name, defaultVal)
	if !found {
		return defaultVal
	}

	val, ok := valI.(string)
	if !ok {
		return fmt.Sprint(valI)
	}

	return val
}

// ReadRequiredString same as ReadRequired but returns string or error
func (p *ParameterBag) ReadRequiredString(name string) (string, error) {
	valI, err := p.ReadRequired(name)
	if err != nil {
		return "", err
	}

	val, ok := valI.(string)
	if !ok {
		val = fmt.Sprint(valI)
	}

	if val == "" {
		return val, fmt.Errorf("required option %s is empty", name)
	}
	return val, nil
}

// ReadStrings same as Read but returns []string
func (p *ParameterBag) ReadStrings(name string, defaultVal ...string) []string {
	valI, found := p.Read(name, defaultVal)
	if !found {
		return defaultVal
	}

	val, ok := valI.([]string)
	if !ok {
		valStr, okStr := valI.(string)
		if okStr {
			return []string{valStr}
		}
		return defaultVal
	}

	return val
}

// ReadRequiredStrings same as Read but returns []string
func (p *ParameterBag) ReadRequiredStrings(name string) ([]string, error) {
	valI, err := p.ReadRequired(name)
	if err != nil {
		return []string{}, err
	}

	val, ok := valI.([]string)
	if !ok {
		valStr, okStr := valI.(string)
		if okStr {
			return []string{valStr}, nil
		}

		return []string{}, fmt.Errorf("cannot convert value %v to []string", valI)
	}

	return val, nil
}

// ReadInt same as Read but returns a int
func (p *ParameterBag) ReadInt(name string, defaultVal int) int {
	valI, found := p.Read(name, defaultVal)
	if !found {
		return defaultVal
	}

	val, ok := valI.(int)
	if !ok {
		intVal, err := strconv.Atoi(fmt.Sprint(valI))
		if err == nil {
			return intVal
		}
		return defaultVal
	}

	return val
}

// ReadRequiredInt same as Read but returns a int and fails if value is missing
func (p *ParameterBag) ReadRequiredInt(name string) (int, error) {
	valI, err := p.ReadRequired(name)
	if err != nil {
		return 0, err
	}

	val, ok := valI.(int)
	if !ok {
		valStr, ok := valI.(string)
		if !ok {
			return 0, fmt.Errorf("cannot convert %v to int", valI)
		}
		intVal, err := strconv.Atoi(valStr)
		if err != nil {
			return 0, fmt.Errorf("cannot convert %v to int", valI)
		}
		return intVal, nil
	}

	return val, nil
}

// ReadInt same as Read but returns a int
func (p *ParameterBag) ReadInt64(name string, defaultVal int64) int64 {
	valI, found := p.Read(name, defaultVal)
	if !found {
		return defaultVal
	}

	val, ok := valI.(int64)
	if !ok {
		uintVal, err := strconv.ParseInt(fmt.Sprint(valI), 10, 64)
		if err == nil {
			return uintVal
		}
		return defaultVal
	}

	return val
}

// ReadRequiredInt same as Read but returns a int64 and fails if value is missing
func (p *ParameterBag) ReadRequiredInt64(name string) (int64, error) {
	valI, err := p.ReadRequired(name)
	if err != nil {
		return 0, err
	}

	val, ok := valI.(int64)
	if !ok {
		intVal, err := strconv.ParseInt(fmt.Sprint(valI), 10, 64)
		if err == nil {
			return intVal, nil
		}
		return 0, fmt.Errorf("cannot convert %v to int64", valI)
	}

	return val, nil
}

// ReadDuration reads int value and converts it to duration identified by the unit, if not set, will return defaultVal
func (p *ParameterBag) ReadDuration(name string, unit time.Duration, defaultVal uint) time.Duration {
	val, err := p.ReadRequiredDuration(name, unit)
	if err != nil {
		return unit * time.Duration(defaultVal)
	}

	return val
}

// ReadRequiredDuration reads int value and converts it to duration identified by the unit, if not set, will return error
func (p *ParameterBag) ReadRequiredDuration(name string, unit time.Duration) (time.Duration, error) {
	valI, err := p.ReadRequired(name)
	if err != nil {
		return 0, err
	}

	val, ok := valI.(time.Duration)
	if ok {
		return val, nil
	}

	valUint, err := p.ReadRequiredUint(name)
	if err != nil {
		return 0, err
	}

	return unit * time.Duration(valUint), nil
}

// ReadBool same as Read but returns a bool or defaultVal
func (p *ParameterBag) ReadBool(name string, defaultVal bool) bool {
	valI, found := p.Read(name, defaultVal)
	if !found {
		return defaultVal
	}

	val, ok := valI.(bool)
	if !ok {
		valStr := strings.TrimSpace(fmt.Sprint(valI))
		switch valStr {
		case "":
			return false
		case "false":
			return false
		case "0":
			return false
		case "[]":
			return false
		}
		return true
	}

	return val
}

// ReadRequiredBool same as ReadRequired but returns bool or error
func (p *ParameterBag) ReadRequiredBool(name string) (bool, error) {
	valI, err := p.ReadRequired(name)
	if err != nil {
		return false, err
	}

	val, ok := valI.(bool)
	if !ok {
		return p.ReadBool(name, false), nil
	}

	return val, nil
}

// ReadUint same as Read but returns a uint
func (p *ParameterBag) ReadUint(name string, defaultVal uint) uint {
	valI, found := p.Read(name, defaultVal)
	if !found {
		return defaultVal
	}

	val, ok := valI.(uint)
	if !ok {
		uintVal, err := strconv.ParseUint(fmt.Sprint(valI), 10, 32)
		if err == nil {
			return uint(uintVal)
		}
		return defaultVal
	}

	return val
}

// ReadRequiredUint same as ReadRequired but returns uint or error
func (p *ParameterBag) ReadRequiredUint(name string) (uint, error) {
	valI, err := p.ReadRequired(name)
	if err != nil {
		return 0, err
	}

	val, ok := valI.(uint)
	if !ok {
		intVal, err := strconv.ParseUint(fmt.Sprint(valI), 10, 32)
		if err == nil {
			return uint(intVal), nil
		}
		return 0, fmt.Errorf("cannot convert %v to uint", valI)
	}

	return val, nil
}

// CheckRequiredValues checks if all required values are not empty
func (p *ParameterBag) CheckRequiredValues(keys []string) error {
	errs := errs2.NewErrorContainer()
	for _, key := range keys {
		_, err := p.ReadRequired(key)
		if err != nil {
			errs.AddError(err)
		}
	}

	return errs.Result(" ")
}

func (p *ParameterBag) MergeParameterBag(m *ParameterBag) {
	valuesProviderComposite := NewValuesProviderComposite(p.BaseValuesProvider, m.BaseValuesProvider)
	p.BaseValuesProvider = valuesProviderComposite
}
