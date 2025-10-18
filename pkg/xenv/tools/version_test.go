package tools

import (
	"testing"
)

func TestParseVersionSpec(t *testing.T) {
	testCases := []struct {
		input    string
		expected *VersionSpec
		hasError bool
	}{
		{
			input: "go:1.21.5",
			expected: &VersionSpec{
				Name:    "go",
				Version: "1.21.5",
			},
			hasError: false,
		},
		{
			input: "node:18",
			expected: &VersionSpec{
				Name:    "node",
				Version: "18",
			},
			hasError: false,
		},
		{
			input: "java:lts",
			expected: &VersionSpec{
				Name:    "java",
				Version: "lts",
			},
			hasError: false,
		},
		{
			input:    "",
			expected: nil,
			hasError: true,
		},
		{
			input:    "go",
			expected: nil,
			hasError: true,
		},
		{
			input:    "go:",
			expected: nil,
			hasError: true,
		},
		{
			input:    ":1.21",
			expected: nil,
			hasError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			result, err := ParseVersionSpec(tc.input)

			if tc.hasError {
				if err == nil {
					t.Errorf("Expected error for input %q, but got none", tc.input)
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error for input %q: %v", tc.input, err)
				return
			}

			if result.Name != tc.expected.Name {
				t.Errorf("Expected Name %q, got %q", tc.expected.Name, result.Name)
			}

			if result.Version != tc.expected.Version {
				t.Errorf("Expected version %q, got %q", tc.expected.Version, result.Version)
			}
		})
	}
}

func TestParseMultipleVersionSpecs(t *testing.T) {
	testCases := []struct {
		name     string
		input    []string
		expected []*VersionSpec
		hasError bool
	}{
		{
			name:  "multiple valid specs",
			input: []string{"go:1.21", "node:18", "java:11"},
			expected: []*VersionSpec{
				{Name: "go", Version: "1.21"},
				{Name: "node", Version: "18"},
				{Name: "java", Version: "11"},
			},
			hasError: false,
		},
		{
			name:     "empty input",
			input:    []string{},
			expected: nil,
			hasError: false,
		},
		{
			name:     "invalid spec in middle",
			input:    []string{"go:1.21", "invalid", "java:11"},
			expected: nil,
			hasError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := ParseMultipleVersionSpecs(tc.input)

			if tc.hasError {
				if err == nil {
					t.Errorf("Expected error for input %v, but got none", tc.input)
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error for input %v: %v", tc.input, err)
				return
			}

			if len(result) != len(tc.expected) {
				t.Errorf("Expected %d specs, got %d", len(tc.expected), len(result))
				return
			}

			for i, expected := range tc.expected {
				if result[i].Name != expected.Name {
					t.Errorf("Expected Name %q at index %d, got %q", expected.Name, i, result[i].Name)
				}
				if result[i].Version != expected.Version {
					t.Errorf("Expected version %q at index %d, got %q", expected.Version, i, result[i].Version)
				}
			}
		})
	}
}

func TestIsValidSDKName(t *testing.T) {
	testCases := []struct {
		input    string
		expected bool
	}{
		{"go", true},
		{"node", true},
		{"java", true},
		{"flutter", true},
		{"my-sdk", true},
		{"my_sdk", true},
		{"sdk123", true},
		{"", false},
		{"sdk with space", false},
		{"sdk@special", false},
		{"sdk.dot", false},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			result := IsValidSDKName(tc.input)
			if result != tc.expected {
				t.Errorf("IsValidSDKName(%q) = %v, expected %v", tc.input, result, tc.expected)
			}
		})
	}
}

func TestIsValidVersion(t *testing.T) {
	testCases := []struct {
		input    string
		expected bool
	}{
		{"1.21.5", true},
		{"18", true},
		{"lts", true},
		{"latest", true},
		{"stable", true},
		{"auto", true},
		{"1.21.5-alpha", true},
		{"1.21.5+build", true},
		{"", false},
		{"version with space", false},
		{"version\twith\ttab", false},
		{"version\nwith\nnewline", false},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			result := IsValidVersion(tc.input)
			if result != tc.expected {
				t.Errorf("IsValidVersion(%q) = %v, expected %v", tc.input, result, tc.expected)
			}
		})
	}
}

func TestNormalizeVersion(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{"  1.21.5  ", "1.21.5"},
		{"LTS", "lts"},
		{"Latest", "latest"},
		{"STABLE", "stable"},
		{"AUTO", "auto"},
		{"1.21.5-alpha", "1.21.5-alpha"},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			result := NormalizeVersion(tc.input)
			if result != tc.expected {
				t.Errorf("NormalizeVersion(%q) = %q, expected %q", tc.input, result, tc.expected)
			}
		})
	}
}

func TestCompareVersions(t *testing.T) {
	testCases := []struct {
		v1       string
		v2       string
		expected int
	}{
		{"1.21.5", "1.21.5", 0},
		{"1.21.4", "1.21.5", -1},
		{"1.21.6", "1.21.5", 1},
		{"go", "node", -1},
		{"node", "go", 1},
	}

	for _, tc := range testCases {
		t.Run(tc.v1+"_vs_"+tc.v2, func(t *testing.T) {
			result := CompareVersions(tc.v1, tc.v2)
			if result != tc.expected {
				t.Errorf("CompareVersions(%q, %q) = %d, expected %d", tc.v1, tc.v2, result, tc.expected)
			}
		})
	}
}

func TestVersionSpecString(t *testing.T) {
	spec := &VersionSpec{
		Name:    "go",
		Version: "1.21.5",
	}

	expected := "go:1.21.5"
	result := spec.String()

	if result != expected {
		t.Errorf("VersionSpec.String() = %q, expected %q", result, expected)
	}
}
