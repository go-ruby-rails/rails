// Copyright (c) the go-ruby-rails/rails authors
//
// SPDX-License-Identifier: BSD-3-Clause

package rails

import (
	"os"
	"reflect"
	"testing"
)

// fakeApp is a minimal App the delegating-accessor tests forward to.
type fakeApp struct {
	root        string
	config      any
	autoloaders any
	publicPath  string
	cache       any
}

func (f fakeApp) Root() string       { return f.root }
func (f fakeApp) Config() any        { return f.config }
func (f fakeApp) Autoloaders() any   { return f.autoloaders }
func (f fakeApp) PublicPath() string { return f.publicPath }
func (f fakeApp) Cache() any         { return f.cache }

// reset returns the package-level module state to its zero values so each test
// starts from a clean, deterministic module. It also installs a getenv that
// reads from env, defaulting to the process environment.
func reset(t *testing.T, env map[string]string) {
	t.Helper()
	mu.Lock()
	application = nil
	envInquirer = nil
	logger = nil
	cacheStore = nil
	cacheStoreSet = false
	backtraceCleaner = nil
	errorReporter = nil
	getenv = func(k string) string { return env[k] }
	mu.Unlock()
	t.Cleanup(func() {
		mu.Lock()
		application = nil
		envInquirer = nil
		logger = nil
		cacheStore = nil
		cacheStoreSet = false
		backtraceCleaner = nil
		errorReporter = nil
		getenv = os.Getenv
		mu.Unlock()
	})
}

func TestApplicationAccessor(t *testing.T) {
	reset(t, nil)
	if Application() != nil {
		t.Error("Application() before boot should be nil")
	}
	app := fakeApp{root: "/srv/app"}
	SetApplication(app)
	if Application() != app {
		t.Errorf("Application() = %v, want %v", Application(), app)
	}
	SetApplication(nil)
	if Application() != nil {
		t.Error("Application() after clear should be nil")
	}
}

func TestEnvDefault(t *testing.T) {
	reset(t, nil)
	e := Env()
	if e.String() != "development" || !e.Is("development") {
		t.Errorf("default Env() = %q, want development", e.String())
	}
	// Second call is memoized: still development even if getenv changes.
	mu.Lock()
	getenv = func(string) string { return "production" }
	mu.Unlock()
	if Env().String() != "development" {
		t.Error("Env() is not memoized")
	}
}

func TestEnvRailsEnv(t *testing.T) {
	reset(t, map[string]string{"RAILS_ENV": "production"})
	if got := Env().String(); got != "production" {
		t.Errorf("Env() = %q, want production", got)
	}
}

func TestEnvRackEnvFallback(t *testing.T) {
	reset(t, map[string]string{"RACK_ENV": "staging"})
	if got := Env().String(); got != "staging" {
		t.Errorf("Env() = %q, want staging", got)
	}
}

func TestSetEnv(t *testing.T) {
	reset(t, nil)
	if got := SetEnv("test"); got.String() != "test" || !got.Local() {
		t.Errorf("SetEnv(test) = %+v, want test/local", got)
	}
	if Env().String() != "test" {
		t.Error("SetEnv did not stick")
	}
}

func TestRoot(t *testing.T) {
	reset(t, nil)
	if Root() != "" {
		t.Error("Root() without app should be empty")
	}
	SetApplication(fakeApp{root: "/srv/app"})
	if Root() != "/srv/app" {
		t.Errorf("Root() = %q, want /srv/app", Root())
	}
}

func TestConfiguration(t *testing.T) {
	reset(t, nil)
	if Configuration() != nil {
		t.Error("Configuration() without app should be nil")
	}
	cfg := struct{ Name string }{"config"}
	SetApplication(fakeApp{config: cfg})
	if Configuration() != cfg {
		t.Errorf("Configuration() = %v, want %v", Configuration(), cfg)
	}
}

func TestAutoloaders(t *testing.T) {
	reset(t, nil)
	if Autoloaders() != nil {
		t.Error("Autoloaders() without app should be nil")
	}
	al := []string{"main", "once"}
	SetApplication(fakeApp{autoloaders: al})
	if !reflect.DeepEqual(Autoloaders(), al) {
		t.Errorf("Autoloaders() = %v, want %v", Autoloaders(), al)
	}
}

func TestPublicPath(t *testing.T) {
	reset(t, nil)
	if PublicPath() != "" {
		t.Error("PublicPath() without app should be empty")
	}
	SetApplication(fakeApp{publicPath: "/srv/app/public"})
	if PublicPath() != "/srv/app/public" {
		t.Errorf("PublicPath() = %q", PublicPath())
	}
}

func TestLogger(t *testing.T) {
	reset(t, nil)
	if Logger() != nil {
		t.Error("Logger() default should be nil")
	}
	l := "a-logger"
	SetLogger(l)
	if Logger() != l {
		t.Errorf("Logger() = %v, want %v", Logger(), l)
	}
}

func TestCache(t *testing.T) {
	reset(t, nil)
	if Cache() != nil {
		t.Error("Cache() without app or override should be nil")
	}
	// Falls back to the application's configured store.
	SetApplication(fakeApp{cache: "app-store"})
	if Cache() != "app-store" {
		t.Errorf("Cache() = %v, want app-store", Cache())
	}
	// An explicit override wins over the application's store.
	SetCache("override-store")
	if Cache() != "override-store" {
		t.Errorf("Cache() = %v, want override-store", Cache())
	}
}

func TestBacktraceCleanerMemoized(t *testing.T) {
	reset(t, nil)
	first := BacktraceCleaner()
	if first == nil {
		t.Fatal("BacktraceCleaner() returned nil")
	}
	if BacktraceCleaner() != first {
		t.Error("BacktraceCleaner() is not memoized")
	}
}

func TestErrorMemoized(t *testing.T) {
	reset(t, nil)
	first := Error()
	if first == nil {
		t.Fatal("Error() returned nil")
	}
	if Error() != first {
		t.Error("Error() is not memoized")
	}
}

func TestGroupsDefault(t *testing.T) {
	reset(t, map[string]string{"RAILS_ENV": "production"})
	got := Groups()
	want := []string{"default", "production"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Groups() = %v, want %v", got, want)
	}
}

func TestGroupsExtrasAndRailsGroupsAndDedup(t *testing.T) {
	reset(t, map[string]string{"RAILS_ENV": "development", "RAILS_GROUPS": "assets,,default"})
	// "assets" positional dup with RAILS_GROUPS, empty entry dropped, and the
	// "default"/env already-present entries de-duplicated.
	got := Groups("assets", "development")
	want := []string{"default", "development", "assets"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Groups() = %v, want %v", got, want)
	}
}

func TestGroupsWithConditional(t *testing.T) {
	reset(t, map[string]string{"RAILS_ENV": "test"})
	got := GroupsWith(map[string][]string{
		"assets": {"development", "test"},
		"tools":  {"development"}, // does not include test
		"perf":   {"test"},
	})
	// Matching keys (assets, perf) are appended in sorted order; tools excluded.
	want := []string{"default", "test", "assets", "perf"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("GroupsWith() = %v, want %v", got, want)
	}
}
