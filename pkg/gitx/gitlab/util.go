package gitlab

import "net/url"

// PathToProjectID convert
func PathToProjectID(pathWithName string) string {
	return url.PathEscape(pathWithName)
}

// BuildProjectID by group and name
func BuildProjectID(group, name string) string {
	return url.PathEscape(group + "/" + name)
}
