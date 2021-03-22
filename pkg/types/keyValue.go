package types

// KeyValueStr just handy implementation for map[string]string
type KeyValueStr struct {
	Key   string
	Value string
}

// KeyValueInt just handy implementation for map[string]int
type KeyValueInt struct {
	Key   string
	Value int
}

// KeyValueInt64 just handy implementation for map[string]int64
type KeyValueInt64 struct {
	Key   string
	Value int64
}

// KeyValueBool just handy implementation for map[string]bool
type KeyValueBool struct {
	Key   string
	Value bool
}

// KeyValueInterface just handy implementation for map[string]interface{} as []KeyValueInterface
type KeyValueInterface struct {
	Key   string
	Value interface{}
}
