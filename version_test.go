// Copyright (c) the go-ruby-rails/rails authors
//
// SPDX-License-Identifier: BSD-3-Clause

package rails

import "testing"

func TestVersionConstant(t *testing.T) {
	if VERSION.Major != 8 || VERSION.Minor != 1 || VERSION.Tiny != 3 || VERSION.Pre != "" {
		t.Fatalf("VERSION = %+v, want 8.1.3 with empty Pre", VERSION)
	}
}

func TestVersionStrings(t *testing.T) {
	const want = "8.1.3"
	if got := VERSION.STRING(); got != want {
		t.Errorf("VERSION.STRING() = %q, want %q", got, want)
	}
	if got := VERSION.String(); got != want {
		t.Errorf("VERSION.String() = %q, want %q", got, want)
	}
	if got := Version(); got != want {
		t.Errorf("Version() = %q, want %q", got, want)
	}
	if got := GemVersion(); got != want {
		t.Errorf("GemVersion() = %q, want %q", got, want)
	}
}

func TestVersionPreRelease(t *testing.T) {
	v := VersionNumber{Major: 9, Minor: 0, Tiny: 0, Pre: "beta1"}
	if got, want := v.STRING(), "9.0.0.beta1"; got != want {
		t.Errorf("STRING() = %q, want %q", got, want)
	}
}
