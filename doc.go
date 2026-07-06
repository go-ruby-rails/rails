// Copyright (c) the go-ruby-rails/rails authors
//
// SPDX-License-Identifier: BSD-3-Clause

// Package rails is the pure-Go (no cgo), MRI-faithful reimplementation of the
// Ruby on Rails meta-gem: the `rails` gem itself.
//
// In Ruby, the `rails` gem ships almost no code of its own. It is the meta-gem
// that ties the framework's components together — it declares a dependency on
// each of them (ActiveSupport, ActiveModel, ActiveJob, ActionPack, ActionView,
// ActionMailer, ActionCable, ActiveStorage, Railties, ...) and provides the
// top-level Rails module: Rails.application, Rails.env, Rails.root, Rails.logger,
// Rails.cache, Rails.configuration, and the Rails::VERSION constant. This package
// mirrors exactly that role.
//
// # The version
//
// [VERSION] is the Rails release this family targets — 8.1.x on MRI 4.0.5, the
// same fidelity basis every sibling component is validated against. [Version]
// and [GemVersion] return its dotted string form, mirroring Rails.version and
// Rails.gem_version.
//
// # The top-level Rails module
//
// The module-level accessors — [Application], [Env], [Root], [Logger], [Cache],
// [Configuration], [BacktraceCleaner], [Autoloaders], [Groups], [PublicPath],
// [Error] — mirror the singleton methods on Ruby's Rails module. The ones that
// expose application state ([Root], [Configuration], [Autoloaders], [PublicPath],
// [Cache]) delegate to the [Application] set via [SetApplication], which the
// go-ruby-railties Application object will satisfy; without an application they
// return the same empty/zero result the Ruby methods return before boot.
//
// [Env] returns an [EnvironmentInquirer] — the pure-Go equivalent of
// ActiveSupport::EnvironmentInquirer — so Env().Is("production") mirrors
// Rails.env.production? and Env().Local() mirrors Rails.env.local?.
//
// # The component manifest
//
// [Components] is the machine-readable catalog of the framework's components:
// each [Component] carries its gem name and derives its go-ruby-* organization,
// import path, and site URLs. The [github.com/go-ruby-rails/rails/all]
// sub-package blank-imports every shipped component, so
//
//	go get github.com/go-ruby-rails/rails
//	import _ "github.com/go-ruby-rails/rails/all"
//
// pulls the whole pure-Go Rails stack into a consumer's build graph behind one
// dependency, and building it cross-compiles the entire framework on every
// supported architecture.
package rails
