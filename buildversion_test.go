// Copyright (c) the go-ruby-rails/rails authors
//
// SPDX-License-Identifier: BSD-3-Clause

package rails

import (
	"runtime/debug"
	"testing"
)

func TestBuildVersionNoBuildInfo(t *testing.T) {
	orig := readBuildInfo
	t.Cleanup(func() { readBuildInfo = orig })
	readBuildInfo = func() (*debug.BuildInfo, bool) { return nil, false }
	if v, ok := BuildVersion("github.com/go-ruby-activesupport/activesupport"); ok || v != "" {
		t.Errorf("BuildVersion with no build info = (%q, %v), want (\"\", false)", v, ok)
	}
}

func TestComponentVersionFound(t *testing.T) {
	orig := readBuildInfo
	t.Cleanup(func() { readBuildInfo = orig })
	readBuildInfo = func() (*debug.BuildInfo, bool) {
		return &debug.BuildInfo{Deps: []*debug.Module{
			{Path: "github.com/go-ruby-activesupport/activesupport", Version: "v1.2.3"},
		}}, true
	}
	c, _ := Lookup("activesupport")
	if v, ok := ComponentVersion(c); !ok || v != "v1.2.3" {
		t.Errorf("ComponentVersion = (%q, %v), want (v1.2.3, true)", v, ok)
	}
	absent, _ := Lookup("activemodel")
	if v, ok := ComponentVersion(absent); ok || v != "" {
		t.Errorf("ComponentVersion(absent) = (%q, %v), want (\"\", false)", v, ok)
	}
}

func TestVersionOfReplace(t *testing.T) {
	info := &debug.BuildInfo{Deps: []*debug.Module{
		{Path: "github.com/go-ruby-a/a", Version: "v0.1.0"},
		{
			Path:    "github.com/go-ruby-b/b",
			Version: "v0.2.0",
			Replace: &debug.Module{Path: "example.com/fork", Version: "v9.9.9"},
		},
	}}
	if v, ok := versionOf(info, "github.com/go-ruby-a/a"); !ok || v != "v0.1.0" {
		t.Errorf("versionOf(a) = (%q, %v), want (v0.1.0, true)", v, ok)
	}
	if v, ok := versionOf(info, "github.com/go-ruby-b/b"); !ok || v != "v9.9.9" {
		t.Errorf("versionOf(b, replaced) = (%q, %v), want (v9.9.9, true)", v, ok)
	}
	if v, ok := versionOf(info, "github.com/missing/missing"); ok || v != "" {
		t.Errorf("versionOf(missing) = (%q, %v), want (\"\", false)", v, ok)
	}
}
