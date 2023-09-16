package dbmeta

// Index table index struct
type Index struct {
	Name    string
	Unique  bool
	Columns []string
}

// Table db table struct
type Table struct {
	Name    string
	Comment string
	Columns []*Column
	Indexes []*Index
}
