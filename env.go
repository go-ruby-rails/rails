// Copyright (c) the go-ruby-rails/rails authors
//
// SPDX-License-Identifier: BSD-3-Clause

package rails

// StringInquirer is the pure-Go equivalent of ActiveSupport::StringInquirer: a
// string wrapper whose predicate queries compare against the wrapped value, so
// that in Ruby `env.production?` is true exactly when the string equals
// "production". Here that predicate is spelled [StringInquirer.Is]; the rbgo
// binding maps the dynamic `name?` methods onto it.
//
// The value is immutable, mirroring the frozen string Rails wraps.
type StringInquirer struct {
	value string
}

// NewStringInquirer wraps value in a StringInquirer.
func NewStringInquirer(value string) StringInquirer { return StringInquirer{value: value} }

// String returns the wrapped string, so a StringInquirer is usable anywhere the
// plain value is (mirroring StringInquirer < String).
func (si StringInquirer) String() string { return si.value }

// Is reports whether the wrapped value equals name, mirroring the dynamic
// `name?` predicate (`env.production?` == `env == "production"`).
func (si StringInquirer) Is(name string) bool { return si.value == name }

// EnvironmentInquirer is the pure-Go equivalent of
// ActiveSupport::EnvironmentInquirer, the StringInquirer subclass Rails uses for
// Rails.env. It adds the environment-specific [EnvironmentInquirer.Local]
// predicate on top of the generic [StringInquirer.Is].
type EnvironmentInquirer struct {
	StringInquirer
}

// NewEnvironmentInquirer wraps name in an EnvironmentInquirer.
func NewEnvironmentInquirer(name string) EnvironmentInquirer {
	return EnvironmentInquirer{StringInquirer: NewStringInquirer(name)}
}

// Local reports whether the environment is one of the built-in local
// environments — "development" or "test" — mirroring Rails.env.local?.
func (e EnvironmentInquirer) Local() bool {
	return e.value == "development" || e.value == "test"
}
