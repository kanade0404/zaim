package domains

import (
	"strings"
)

// Domain is an interface that all domain models must implement.
type Domain interface {
	Indexes() []*Index
}

// baseIndex is an interface that all index types must implement.
type baseIndex interface {
	Name() string
	Columns() []string
}

// Index represents a non-unique index.
type Index struct {
	name    string
	columns []string
}

var _ baseIndex = (*Index)(nil)

// Name returns the name of the index.
func (i *Index) Name() string {
	return i.name
}

// Columns returns the columns of the index.
func (i *Index) Columns() []string {
	return i.columns
}

// createIndex creates a new Index.
func createIndex(tableName string, columns []string) *Index {
	return &Index{
		name:    "idx_" + tableName + "_" + strings.Join(columns, "_"),
		columns: columns,
	}
}
