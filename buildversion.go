// Copyright (c) the go-ruby-rails/rails authors
//
// SPDX-License-Identifier: BSD-3-Clause

package rails

import "runtime/debug"

// readBuildInfo is indirected so tests can exercise the missing-build-info path.
var readBuildInfo = debug.ReadBuildInfo

// ComponentVersion returns the version a framework component's module was pinned
// to in the build graph of the running binary, if that module is present.
//
// The version is only available when the component is actually part of the
// binary's build graph — for example because the binary imports the rails/all
// aggregate, or imports the component directly. When it is not present, or when
// the binary was built without module information, ok is false.
func ComponentVersion(c Component) (version string, ok bool) {
	return BuildVersion(c.ImportPath())
}

// BuildVersion is [ComponentVersion] keyed by import path.
func BuildVersion(importPath string) (version string, ok bool) {
	info, present := readBuildInfo()
	if !present {
		return "", false
	}
	return versionOf(info, importPath)
}

// versionOf finds importPath among the build info's dependencies, honouring a
// replace directive if one is in effect.
func versionOf(info *debug.BuildInfo, importPath string) (string, bool) {
	for _, dep := range info.Deps {
		if dep.Path != importPath {
			continue
		}
		if dep.Replace != nil {
			return dep.Replace.Version, true
		}
		return dep.Version, true
	}
	return "", false
}
