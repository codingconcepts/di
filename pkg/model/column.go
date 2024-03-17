package model

// ColumnTypes is a map of column name to column information.
type ColumnTypes map[string]*Column

// Column represents the properties of a table column.
type Column struct {
	Ordinal  int
	Name     string
	Type     string
	Nullable bool
	Format   string
}

// SetFormat sets a custom formatter for a column.
func (c *Column) SetFormat(fmt string) {
	c.Format = fmt
}
