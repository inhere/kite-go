package textparse

type TextParser struct {
	// Delimiter for split items
	Delimiter string
}

type ItemParser interface {
	// Parse item text to values
	Parse(item string) []string
}

type LineParser interface {
	// Parse line to values
	Parse(item string) []string
}

// TextItem text item struct
type TextItem struct {
	// start line of text item
	Line int
	// Index of items
	Index int
	// Item content
	Content string
	// parsed values of item text
	Values []string
}
