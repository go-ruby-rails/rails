// Copyright (c) the go-ruby-rails/rails authors
//
// SPDX-License-Identifier: BSD-3-Clause

package rails

import "strconv"

// VersionNumber is the structured Rails release number, the pure-Go equivalent
// of Ruby's Rails::VERSION module. Major, Minor, and Tiny mirror the MAJOR,
// MINOR, and TINY constants; Pre mirrors PRE (empty when the release is final).
type VersionNumber struct {
	Major int
	Minor int
	Tiny  int
	Pre   string
}

// VERSION is the Rails release this family targets: 8.1.x on MRI 4.0.5, the
// fidelity basis shared by every go-ruby-* framework component.
var VERSION = VersionNumber{Major: 8, Minor: 1, Tiny: 3, Pre: ""}

// STRING returns the dotted release string, mirroring Rails::VERSION::STRING,
// which is [MAJOR, MINOR, TINY, PRE].compact.join("."). The pre-release segment
// is appended only when present.
func (v VersionNumber) STRING() string {
	s := strconv.Itoa(v.Major) + "." + strconv.Itoa(v.Minor) + "." + strconv.Itoa(v.Tiny)
	if v.Pre != "" {
		s += "." + v.Pre
	}
	return s
}

// String implements fmt.Stringer with the same dotted form as [VersionNumber.STRING].
func (v VersionNumber) String() string { return v.STRING() }

// Version returns the Rails release as a dotted string, mirroring Rails.version
// (which returns a Gem::Version built from Rails::VERSION::STRING).
func Version() string { return VERSION.STRING() }

// GemVersion returns the Rails release as a dotted string, mirroring
// Rails.gem_version.
func GemVersion() string { return VERSION.STRING() }
