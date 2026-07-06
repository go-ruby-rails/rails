// Copyright (c) the go-ruby-rails/rails authors
//
// SPDX-License-Identifier: BSD-3-Clause

package rails

import (
	"sort"
	"strings"
	"testing"
)

func TestCatalogIntegrity(t *testing.T) {
	if Count() != len(catalog) {
		t.Fatalf("Count() = %d, want %d", Count(), len(catalog))
	}
	if Count() == 0 {
		t.Fatal("catalog is empty")
	}
	seen := map[string]bool{}
	for _, c := range catalog {
		switch {
		case c.Name == "":
			t.Errorf("component has empty Name: %+v", c)
		case c.Gem == "":
			t.Errorf("component %q has empty Gem", c.Name)
		case c.Description == "":
			t.Errorf("component %q has empty Description", c.Name)
		}
		if seen[c.Name] {
			t.Errorf("duplicate component Name %q", c.Name)
		}
		seen[c.Name] = true
	}
}

func TestComponentsSortedCopy(t *testing.T) {
	got := Components()
	if len(got) != Count() {
		t.Fatalf("Components() len = %d, want %d", len(got), Count())
	}
	if !sort.SliceIsSorted(got, func(i, j int) bool { return got[i].Name < got[j].Name }) {
		t.Error("Components() is not sorted by Name")
	}
	got[0] = Component{Name: "tampered"}
	if _, ok := Lookup("tampered"); ok {
		t.Error("mutating Components() result leaked into the manifest")
	}
}

func TestAvailableComponents(t *testing.T) {
	avail := AvailableComponents()
	if len(avail) == 0 {
		t.Fatal("no available components")
	}
	if !sort.SliceIsSorted(avail, func(i, j int) bool { return avail[i].Name < avail[j].Name }) {
		t.Error("AvailableComponents() not sorted")
	}
	for _, c := range avail {
		if !c.Available {
			t.Errorf("AvailableComponents() returned unavailable %q", c.Name)
		}
	}
	// Every catalogued component (including actionmailer and railties) has now
	// shipped, so the available set equals the full manifest.
	if len(avail) != Count() {
		t.Errorf("expected all components available; avail=%d total=%d", len(avail), Count())
	}
}

func TestLookup(t *testing.T) {
	c, ok := Lookup("activesupport")
	if !ok {
		t.Fatal("Lookup(activesupport) not found")
	}
	if c.Gem != "ActiveSupport" || !c.Available {
		t.Errorf("Lookup(activesupport) = %+v, unexpected", c)
	}
	if r, ok := Lookup("railties"); !ok || !r.Available {
		t.Errorf("Lookup(railties) = (%+v, %v), want present and available", r, ok)
	}
	if _, ok := Lookup("nonesuch"); ok {
		t.Error("Lookup(nonesuch) unexpectedly found")
	}
}

func TestImportPathsUniqueSorted(t *testing.T) {
	paths := ImportPaths()
	if len(paths) != Count() {
		t.Fatalf("ImportPaths() len = %d, want %d", len(paths), Count())
	}
	if !sort.StringsAreSorted(paths) {
		t.Error("ImportPaths() not sorted")
	}
	seen := map[string]bool{}
	for _, p := range paths {
		if seen[p] {
			t.Errorf("duplicate import path %q", p)
		}
		seen[p] = true
		if !strings.HasPrefix(p, "github.com/go-ruby-") {
			t.Errorf("import path %q has unexpected prefix", p)
		}
	}
}

func TestAggregateImportPathsMatchesAvailable(t *testing.T) {
	agg := AggregateImportPaths()
	if !sort.StringsAreSorted(agg) {
		t.Error("AggregateImportPaths() not sorted")
	}
	if len(agg) != len(AvailableComponents()) {
		t.Fatalf("AggregateImportPaths() len = %d, want %d", len(agg), len(AvailableComponents()))
	}
	for i, c := range AvailableComponents() {
		if agg[i] != c.ImportPath() {
			t.Errorf("agg[%d] = %q, want %q", i, agg[i], c.ImportPath())
		}
	}
}

func TestComponentDerivedURLs(t *testing.T) {
	c := Component{Name: "activesupport"}
	cases := map[string]string{
		"Org":        c.Org(),
		"ImportPath": c.ImportPath(),
		"RepoURL":    c.RepoURL(),
		"LandingURL": c.LandingURL(),
		"DocsURL":    c.DocsURL(),
	}
	want := map[string]string{
		"Org":        "go-ruby-activesupport",
		"ImportPath": "github.com/go-ruby-activesupport/activesupport",
		"RepoURL":    "https://github.com/go-ruby-activesupport/activesupport",
		"LandingURL": "https://go-ruby-activesupport.github.io/",
		"DocsURL":    "https://go-ruby-activesupport.github.io/docs/",
	}
	for k, got := range cases {
		if got != want[k] {
			t.Errorf("%s() = %q, want %q", k, got, want[k])
		}
	}
}
