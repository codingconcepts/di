package flags

// StringSlice can be used as a flag argument on the command line.
type StringSlice []string

func (i *StringSlice) String() string {
	return "my string representation"
}

func (i *StringSlice) Set(value string) error {
	*i = append(*i, value)
	return nil
}
