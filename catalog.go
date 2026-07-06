// Copyright (c) the go-ruby-rails/rails authors
//
// SPDX-License-Identifier: BSD-3-Clause

package rails

import "sort"

// Component describes one component of the Rails framework: a pure-Go,
// MRI-faithful reimplementation of one of the gems the `rails` meta-gem ties
// together (ActiveSupport, ActionPack, ...).
//
// The organization, import path, and site URLs are not stored — they follow the
// go-ruby-* family's uniform convention and are derived from Name, so the
// manifest cannot drift out of sync with them.
type Component struct {
	// Name is the gem / repository name, for example "activesupport". It is both
	// the go-ruby-* org suffix and the Go module name.
	Name string
	// Gem is the Ruby constant the component provides, for example
	// "ActiveSupport".
	Gem string
	// Description is a one-line summary of the component's responsibility.
	Description string
	// Available reports whether the go-ruby-* component has shipped and is
	// therefore blank-imported by the rails/all aggregate. Components whose
	// go-ruby-* repository is still empty are catalogued with Available=false.
	Available bool
}

// Org returns the GitHub organization that owns the component, for example
// "go-ruby-activesupport".
func (c Component) Org() string { return "go-ruby-" + c.Name }

// ImportPath returns the Go module / import path, for example
// "github.com/go-ruby-activesupport/activesupport".
func (c Component) ImportPath() string { return "github.com/" + c.Org() + "/" + c.Name }

// RepoURL returns the canonical source-repository URL.
func (c Component) RepoURL() string { return "https://github.com/" + c.Org() + "/" + c.Name }

// LandingURL returns the component's landing-site URL.
func (c Component) LandingURL() string { return "https://" + c.Org() + ".github.io/" }

// DocsURL returns the component's documentation-site URL.
func (c Component) DocsURL() string { return "https://" + c.Org() + ".github.io/docs/" }

// catalog is the authoritative manifest of the Rails framework components the
// meta-gem ties together, in canonical Rails load order. Public accessors
// return sorted copies so callers can never mutate it.
var catalog = []Component{
	{Name: "activesupport", Gem: "ActiveSupport", Available: true, Description: "core Ruby extensions and framework utilities"},
	{Name: "activemodel", Gem: "ActiveModel", Available: true, Description: "model interface: validations, naming, errors"},
	{Name: "activejob", Gem: "ActiveJob", Available: true, Description: "background-job framework and queue adapters"},
	{Name: "actionpack", Gem: "ActionPack", Available: true, Description: "controllers, routing, and request dispatch"},
	{Name: "actionview", Gem: "ActionView", Available: true, Description: "view rendering and template helpers"},
	{Name: "actionmailer", Gem: "ActionMailer", Available: true, Description: "email delivery framework"},
	{Name: "actioncable", Gem: "ActionCable", Available: true, Description: "WebSocket / pub-sub integration"},
	{Name: "activestorage", Gem: "ActiveStorage", Available: true, Description: "file attachments and blob storage"},
	{Name: "railties", Gem: "Rails::Railtie", Available: true, Description: "application boot, engines, and the generators/CLI"},
}

// Count reports the number of framework components in the manifest.
func Count() int { return len(catalog) }

// Components returns every framework component in the manifest, sorted by Name.
// The returned slice is a copy; mutating it cannot affect the manifest.
func Components() []Component {
	out := make([]Component, len(catalog))
	copy(out, catalog)
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out
}

// AvailableComponents returns the components that have shipped and are included
// in the rails/all aggregate, sorted by Name.
func AvailableComponents() []Component {
	var out []Component
	for _, c := range catalog {
		if c.Available {
			out = append(out, c)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out
}

// Lookup returns the component with the given Name.
func Lookup(name string) (Component, bool) {
	for _, c := range catalog {
		if c.Name == name {
			return c, true
		}
	}
	return Component{}, false
}

// ImportPaths returns the Go import path of every component in the manifest,
// sorted.
func ImportPaths() []string {
	out := make([]string, len(catalog))
	for i, c := range catalog {
		out[i] = c.ImportPath()
	}
	sort.Strings(out)
	return out
}

// AggregateImportPaths returns the Go import paths blank-imported by the
// rails/all aggregate — the available components only — sorted. A test keeps
// all/all.go exactly in sync with this list.
func AggregateImportPaths() []string {
	var out []string
	for _, c := range catalog {
		if c.Available {
			out = append(out, c.ImportPath())
		}
	}
	sort.Strings(out)
	return out
}
