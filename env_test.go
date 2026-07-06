// Copyright (c) the go-ruby-rails/rails authors
//
// SPDX-License-Identifier: BSD-3-Clause

package rails

import "testing"

func TestStringInquirer(t *testing.T) {
	si := NewStringInquirer("production")
	if si.String() != "production" {
		t.Errorf("String() = %q, want production", si.String())
	}
	if !si.Is("production") {
		t.Error("Is(production) = false, want true")
	}
	if si.Is("development") {
		t.Error("Is(development) = true, want false")
	}
}

func TestEnvironmentInquirerLocal(t *testing.T) {
	cases := map[string]bool{
		"development": true,
		"test":        true,
		"production":  false,
		"staging":     false,
	}
	for name, wantLocal := range cases {
		e := NewEnvironmentInquirer(name)
		if got := e.Local(); got != wantLocal {
			t.Errorf("NewEnvironmentInquirer(%q).Local() = %v, want %v", name, got, wantLocal)
		}
		if !e.Is(name) {
			t.Errorf("Is(%q) = false, want true", name)
		}
	}
}
