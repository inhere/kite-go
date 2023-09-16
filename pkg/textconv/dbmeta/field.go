package dbmeta

// Column table column struct
type Column struct {
	// Name of the field
	Name string
	// Type of the field
	Type string
	// TypeLen of the field. eg: 11
	TypeLen int
	// TypeExt of the field. eg: UNSIGNED
	TypeExt string
	// Nullable of the field
	Nullable bool
	// Default value of the field
	Default string
	// Comment of the field
	Comment string
}
