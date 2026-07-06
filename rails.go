// Copyright (c) the go-ruby-rails/rails authors
//
// SPDX-License-Identifier: BSD-3-Clause

package rails

import (
	"os"
	"sort"
	"strings"
	"sync"
)

// App is the minimal surface of the Rails application object that the top-level
// Rails module delegates to. In Ruby this is a Rails::Application (a Railtie)
// instance; the meta-gem never constructs one — go-ruby-railties owns the
// concrete type. Modelling it as an interface here keeps the meta-gem decoupled
// from railties while still letting the delegating accessors ([Root],
// [Configuration], [Autoloaders], [PublicPath], [Cache]) forward to it.
//
// The configuration, autoloader set, and cache store are opaque to the meta-gem
// (owned by railties and the components), so they are typed as any — exactly as
// Ruby's Rails.configuration simply returns application.config.
type App interface {
	// Root is the application's root directory (config.root), a filesystem path.
	Root() string
	// Config is the application's configuration object (application.config).
	Config() any
	// Autoloaders is the application's autoloader set (application.autoloaders).
	Autoloaders() any
	// PublicPath is the first entry of paths["public"].
	PublicPath() string
	// Cache is the configured cache store (config.cache_store instance).
	Cache() any
}

// backtraceCleanerObj and errorReporterObj are the minimal singletons the
// meta-gem vends for Rails.backtrace_cleaner and Rails.error. The real objects
// (ActiveSupport::BacktraceCleaner and ActiveSupport::ErrorReporter) are owned
// by go-ruby-activesupport; the accessors return them as opaque values.
type backtraceCleanerObj struct{}

type errorReporterObj struct{}

var (
	mu               sync.RWMutex
	application      App
	envInquirer      *EnvironmentInquirer
	logger           any
	cacheStore       any
	cacheStoreSet    bool
	backtraceCleaner any
	errorReporter    any

	// getenv is indirected so tests can drive the environment-detection paths
	// deterministically.
	getenv = os.Getenv
)

// SetApplication registers the running application, mirroring
// `Rails.application = app`. Passing nil clears it.
func SetApplication(app App) { mu.Lock(); application = app; mu.Unlock() }

// Application returns the registered application, mirroring Rails.application. It
// is nil until [SetApplication] is called (before the app boots).
func Application() App { mu.RLock(); defer mu.RUnlock(); return application }

// Env returns the current environment as an [EnvironmentInquirer], mirroring
// Rails.env. On first use the value is resolved from RAILS_ENV, then RACK_ENV,
// then defaults to "development", and is memoized thereafter (as in Rails).
func Env() EnvironmentInquirer {
	mu.Lock()
	defer mu.Unlock()
	if envInquirer == nil {
		name := getenv("RAILS_ENV")
		if name == "" {
			name = getenv("RACK_ENV")
		}
		if name == "" {
			name = "development"
		}
		e := NewEnvironmentInquirer(name)
		envInquirer = &e
	}
	return *envInquirer
}

// SetEnv overrides the current environment, mirroring `Rails.env = "production"`.
// It returns the new inquirer.
func SetEnv(name string) EnvironmentInquirer {
	mu.Lock()
	defer mu.Unlock()
	e := NewEnvironmentInquirer(name)
	envInquirer = &e
	return e
}

// Root returns the application's root directory, mirroring Rails.root. It is the
// empty string when no application is registered (Ruby returns nil).
func Root() string {
	mu.RLock()
	defer mu.RUnlock()
	if application == nil {
		return ""
	}
	return application.Root()
}

// Configuration returns the application's configuration, mirroring
// Rails.configuration. It is nil when no application is registered.
func Configuration() any {
	mu.RLock()
	defer mu.RUnlock()
	if application == nil {
		return nil
	}
	return application.Config()
}

// Autoloaders returns the application's autoloader set, mirroring
// Rails.autoloaders. It is nil when no application is registered.
func Autoloaders() any {
	mu.RLock()
	defer mu.RUnlock()
	if application == nil {
		return nil
	}
	return application.Autoloaders()
}

// PublicPath returns the application's public directory, mirroring
// Rails.public_path. It is the empty string when no application is registered
// (Ruby returns nil).
func PublicPath() string {
	mu.RLock()
	defer mu.RUnlock()
	if application == nil {
		return ""
	}
	return application.PublicPath()
}

// SetLogger sets the Rails logger, mirroring `Rails.logger = logger`.
func SetLogger(l any) { mu.Lock(); logger = l; mu.Unlock() }

// Logger returns the Rails logger, mirroring Rails.logger. It is nil until one
// is set.
func Logger() any { mu.RLock(); defer mu.RUnlock(); return logger }

// SetCache sets the Rails cache store, mirroring `Rails.cache = store`.
func SetCache(c any) { mu.Lock(); cacheStore, cacheStoreSet = c, true; mu.Unlock() }

// Cache returns the Rails cache store, mirroring Rails.cache. An explicit store
// set with [SetCache] takes precedence; otherwise the application's configured
// store is used; without either it is nil.
func Cache() any {
	mu.RLock()
	defer mu.RUnlock()
	if cacheStoreSet {
		return cacheStore
	}
	if application == nil {
		return nil
	}
	return application.Cache()
}

// BacktraceCleaner returns the shared backtrace cleaner, mirroring
// Rails.backtrace_cleaner. It is created on first use and memoized. The concrete
// object is owned by ActiveSupport, so it is returned as an opaque value.
func BacktraceCleaner() any {
	mu.Lock()
	defer mu.Unlock()
	if backtraceCleaner == nil {
		backtraceCleaner = &backtraceCleanerObj{}
	}
	return backtraceCleaner
}

// Error returns the shared error reporter, mirroring Rails.error (which
// delegates to ActiveSupport.error_reporter). It is created on first use and
// memoized. The concrete object is owned by ActiveSupport, so it is returned as
// an opaque value.
func Error() any {
	mu.Lock()
	defer mu.Unlock()
	if errorReporter == nil {
		errorReporter = &errorReporterObj{}
	}
	return errorReporter
}

// Groups returns the ordered, de-duplicated list of groups active in the
// current environment, mirroring Rails.groups(*groups).
//
// The result is `default`, then the current environment, then any positional
// extras, then the comma-separated RAILS_GROUPS environment variable — with
// empty entries dropped and duplicates removed (keeping first occurrence).
func Groups(extra ...string) []string {
	return GroupsWith(nil, extra...)
}

// GroupsWith is [Groups] with the conditional-hash form of Rails.groups: for
// each key in conditional whose value contains the current environment, the key
// is appended to the group list. This mirrors, for example,
// `Rails.groups(assets: %w(development test))`. Keys are considered in sorted
// order so the result is deterministic.
func GroupsWith(conditional map[string][]string, extra ...string) []string {
	env := Env().String()

	groups := []string{"default", env}
	groups = append(groups, extra...)
	groups = append(groups, strings.Split(getenv("RAILS_GROUPS"), ",")...)
	for _, key := range sortedMatchingKeys(conditional, env) {
		groups = append(groups, key)
	}

	seen := make(map[string]bool, len(groups))
	out := make([]string, 0, len(groups))
	for _, g := range groups {
		if g == "" || seen[g] {
			continue
		}
		seen[g] = true
		out = append(out, g)
	}
	return out
}

// sortedMatchingKeys returns, in ascending order, the keys of conditional whose
// value slice contains env.
func sortedMatchingKeys(conditional map[string][]string, env string) []string {
	var keys []string
	for key, envs := range conditional {
		for _, e := range envs {
			if e == env {
				keys = append(keys, key)
				break
			}
		}
	}
	sort.Strings(keys)
	return keys
}
