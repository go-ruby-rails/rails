// Copyright (c) the go-ruby-rails/rails authors
//
// SPDX-License-Identifier: BSD-3-Clause

package rails

import (
	"go/parser"
	"go/token"
	"sort"
	"strconv"
	"testing"
)

// TestAllPackageInSync guards the invariant that the all sub-package
// blank-imports exactly the available components in the manifest — no more, no
// less. If a component's Available flag changes without updating all/all.go (or
// vice versa), this fails.
func TestAllPackageInSync(t *testing.T) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "all/all.go", nil, parser.ImportsOnly)
	if err != nil {
		t.Fatalf("parsing all/all.go: %v", err)
	}
	var imported []string
	for _, spec := range f.Imports {
		path, err := strconv.Unquote(spec.Path.Value)
		if err != nil {
			t.Fatalf("unquoting import %s: %v", spec.Path.Value, err)
		}
		imported = append(imported, path)
	}
	sort.Strings(imported)

	want := AggregateImportPaths()
	if len(imported) != len(want) {
		t.Fatalf("all/all.go imports %d components, manifest has %d available", len(imported), len(want))
	}
	for i := range want {
		if imported[i] != want[i] {
			t.Errorf("import[%d] = %q, manifest has %q", i, imported[i], want[i])
		}
	}
}
