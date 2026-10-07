package gitcmd

import (
	"testing"

	"github.com/gookit/goutil/testutil/assert"
)

func TestParseBumpLevel(t *testing.T) {
	testCases := []struct {
		input  string
		want   VersionBumpLevel
		wantOk bool
	}{
		{input: "", want: BumpPatch, wantOk: true},
		{input: "patch", want: BumpPatch, wantOk: true},
		{input: "fix", want: BumpPatch, wantOk: true},
		{input: "3", want: BumpPatch, wantOk: true},
		{input: "minor", want: BumpMinor, wantOk: true},
		{input: "min", want: BumpMinor, wantOk: true},
		{input: "2", want: BumpMinor, wantOk: true},
		{input: "MAJOR", want: BumpMajor, wantOk: true},
		{input: "maj", want: BumpMajor, wantOk: true},
		{input: "1", want: BumpMajor, wantOk: true},
		{input: " foo ", wantOk: false},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			got, ok := ParseBumpLevel(tc.input)
			assert.Eq(t, tc.wantOk, ok)
			if tc.wantOk {
				assert.Eq(t, tc.want, got)
			}
		})
	}
}

func TestBumpVersion(t *testing.T) {
	testCases := []struct {
		version string
		level   VersionBumpLevel
		want    string
	}{
		{version: "v1.2.3", level: BumpPatch, want: "1.2.4"},
		{version: "1.2.3", level: BumpMinor, want: "1.3.0"},
		{version: "v1.2.3", level: BumpMajor, want: "2.0.0"},
		{version: "v2.0.0", level: BumpMajor, want: "3.0.0"},
		{version: "", level: BumpPatch, want: "0.0.1"},
		{version: "v1.2", level: BumpPatch, want: "1.2.1"},
		{version: "v1.2", level: BumpMinor, want: "1.3.0"},
		{version: "v1.2.3-beta", level: BumpPatch, want: "1.2.4"},
	}

	for _, tc := range testCases {
		t.Run(tc.version, func(t *testing.T) {
			assert.Eq(t, tc.want, BumpVersion(tc.version, tc.level))
		})
	}
}
