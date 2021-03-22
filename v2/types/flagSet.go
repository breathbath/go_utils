package types

// FlagSet holds the info about certain keys to appear in tests
type FlagSet map[string]bool

// NewFlagSet contructor
func NewFlagSet() FlagSet {
	return FlagSet{}
}

// AllTrue useful for assertions if all expected values appeared e.g. in a stream
func (fs FlagSet) AllTrue() bool {
	if len(fs) == 0 {
		return false
	}

	for _, val := range fs {
		if !val {
			return false
		}
	}

	return true
}

// GetNotMatchedKeys useful for giving failure texts
func (fs FlagSet) GetNotMatchedKeys() []string {
	notMatchedKeys := make([]string, 0, len(fs))
	for key, val := range fs {
		if !val {
			notMatchedKeys = append(notMatchedKeys, key)
		}
	}

	return notMatchedKeys
}
