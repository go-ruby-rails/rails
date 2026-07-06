// Copyright (c) the go-ruby-rails/rails authors
//
// SPDX-License-Identifier: BSD-3-Clause

// Package all blank-imports every shipped component of the Rails framework, so
// that a single
//
//	import _ "github.com/go-ruby-rails/rails/all"
//
// pulls the whole pure-Go, CGO-free, MRI-faithful Rails stack into a consumer's
// build graph behind one dependency.
//
// Importing this package links all of the shipped framework components into the
// binary. A consumer that only needs a few should import those directly; this
// package exists for the "give me the whole framework" case and, in CI, as the
// integration proof that every component compiles together at its pinned
// pseudo-version on every supported architecture.
//
// The set of imports below is kept exactly in sync with
// [github.com/go-ruby-rails/rails.AggregateImportPaths] by a test in the parent
// package; do not edit it by hand without updating the manifest.
package all

import (
	_ "github.com/go-ruby-actioncable/actioncable"
	_ "github.com/go-ruby-actionpack/actionpack"
	_ "github.com/go-ruby-actionview/actionview"
	_ "github.com/go-ruby-activejob/activejob"
	_ "github.com/go-ruby-activemodel/activemodel"
	_ "github.com/go-ruby-activestorage/activestorage"
	_ "github.com/go-ruby-activesupport/activesupport"
)
