package flags

import (
	"strings"
)

// StringSlice can be used as a flag argument on the command line.
type StringSlice []string

func NewStringSlice() *StringSlice {
	return &StringSlice{
		"date=2006-01-02",
		"time=15:04:05",
		"timestamp=2006-01-02T15:04:05",
		"timestamptz=2006-01-02T15:04:05",
	}
}

func (i *StringSlice) String() string {
	return strings.Join(*i, ", ")
}

func (i *StringSlice) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func (i *StringSlice) ToMap() map[string]string {
	m := map[string]string{}

	for _, s := range *i {
		parts := strings.Split(s, "=")
		m[parts[0]] = parts[1]
	}

	return m
}
