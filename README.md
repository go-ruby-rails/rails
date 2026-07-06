<p align="center"><img src="https://raw.githubusercontent.com/go-ruby-rails/brand/main/social/go-ruby-rails-rails.png" alt="go-ruby-rails/rails" width="720"></p>

# rails — go-ruby-rails

[![Docs](https://img.shields.io/badge/docs-mkdocs--material-DC2626)](https://go-ruby-rails.github.io/docs/)
[![License](https://img.shields.io/badge/license-BSD--3--Clause-blue)](LICENSE)
[![Go](https://img.shields.io/badge/go-1.26.4%2B-00ADD8)](https://go.dev/dl/)
[![Coverage](https://img.shields.io/badge/coverage-100%25-1a7f37)](#tests--coverage)

**A pure-Go (no cgo) reimplementation of the Ruby on Rails
[`rails`](https://rubygems.org/gems/rails) meta-gem** — targeting Rails **8.1.x**
on MRI 4.0.5, the same fidelity basis every go-ruby-* framework component is
validated against.

In Ruby, `rails` ships almost no code of its own. It is the **meta-gem**: it
declares a dependency on each framework component (ActiveSupport, ActiveModel,
ActiveJob, ActionPack, ActionView, ActionMailer, ActionCable, ActiveStorage,
Railties) and provides the top-level **`Rails` module** — `Rails.application`,
`Rails.env`, `Rails.root`, `Rails.logger`, `Rails.cache`, `Rails.configuration`,
and the `Rails::VERSION` constant. This package mirrors exactly that role, with
no Ruby runtime.

It is the `rails` backend for
[go-embedded-ruby](https://github.com/go-embedded-ruby/ruby) and a
**standalone, reusable** Go module — the capstone that ties the go-ruby-*
Rails family together.

## What it ties together

The framework components are independent, MRI-faithful, pure-Go modules, each
under its own `github.com/go-ruby-<name>/<name>` organization. This meta-gem
carries a machine-readable manifest of them and a blank-import aggregate:

| Component | Ruby gem | go-ruby-* module | In `rails/all` |
| --- | --- | --- | --- |
| activesupport | `ActiveSupport` | [go-ruby-activesupport/activesupport](https://github.com/go-ruby-activesupport/activesupport) | ✅ |
| activemodel | `ActiveModel` | [go-ruby-activemodel/activemodel](https://github.com/go-ruby-activemodel/activemodel) | ✅ |
| activejob | `ActiveJob` | [go-ruby-activejob/activejob](https://github.com/go-ruby-activejob/activejob) | ✅ |
| actionpack | `ActionPack` | [go-ruby-actionpack/actionpack](https://github.com/go-ruby-actionpack/actionpack) | ✅ |
| actionview | `ActionView` | [go-ruby-actionview/actionview](https://github.com/go-ruby-actionview/actionview) | ✅ |
| actioncable | `ActionCable` | [go-ruby-actioncable/actioncable](https://github.com/go-ruby-actioncable/actioncable) | ✅ |
| activestorage | `ActiveStorage` | [go-ruby-activestorage/activestorage](https://github.com/go-ruby-activestorage/activestorage) | ✅ |
| actionmailer | `ActionMailer` | [go-ruby-actionmailer/actionmailer](https://github.com/go-ruby-actionmailer/actionmailer) | ✅ |
| railties | `Rails::Railtie` | [go-ruby-railties/railties](https://github.com/go-ruby-railties/railties) | ✅ |

The manifest is the single source of truth: `Components()`, `AvailableComponents()`,
`Lookup(name)`, `ImportPaths()`, and `AggregateImportPaths()` derive every
component's org, import path, and site URLs from its name, so the catalog cannot
drift. Every catalogued component has shipped, so all are `Available = true` and
blank-imported by `rails/all`; a test keeps the manifest and the aggregate
exactly in sync.

## The Rails module

```go
import "github.com/go-ruby-rails/rails"

rails.VERSION            // VersionNumber{Major: 8, Minor: 1, Tiny: 3}
rails.Version()          // "8.1.3"   — Rails.version
rails.GemVersion()       // "8.1.3"   — Rails.gem_version

rails.SetEnv("production")
rails.Env().Is("production") // true  — Rails.env.production?
rails.Env().Local()          // false — Rails.env.local? (development || test)

rails.Groups()                            // ["default", "production"]
rails.GroupsWith(map[string][]string{     // Rails.groups(assets: %w[development test])
    "assets": {"development", "test"},
})

rails.SetApplication(app) // app satisfies rails.App (go-ruby-railties provides it)
rails.Root()              // Rails.root          (delegates to app)
rails.Configuration()     // Rails.configuration (delegates to app)
rails.Autoloaders()       // Rails.autoloaders   (delegates to app)
rails.PublicPath()        // Rails.public_path   (delegates to app)
rails.Cache()             // Rails.cache         (override, else app store)
rails.Logger()            // Rails.logger
rails.BacktraceCleaner()  // Rails.backtrace_cleaner
rails.Error()             // Rails.error
```

`Rails.env` returns an **`EnvironmentInquirer`** — the pure-Go equivalent of
`ActiveSupport::EnvironmentInquirer` — so `Env().Is("production")` mirrors
`Rails.env.production?`. The accessors that expose application state delegate to
the `App` set via `SetApplication`; without a booted application they return the
same empty/zero result the Ruby methods return before boot. The concrete
application type is owned by go-ruby-railties — the meta-gem stays decoupled from
it through the small `App` interface.

## The `rails/all` aggregate

```go
import _ "github.com/go-ruby-rails/rails/all"
```

blank-imports every shipped component, so a single `go get
github.com/go-ruby-rails/rails` and one import pull the whole pure-Go Rails stack
into your build graph. Building `rails/all` is also the family's **integration
proof**: it compiles every component together, at its pinned pseudo-version, on
every supported architecture.

`ComponentVersion(c)` reports the exact version each component was pinned to in
the running binary's build graph (read from the embedded build info).

## Tests & coverage

The package holds **100% line coverage** (version, env inquirer, groups, the
module accessors, and the manifest — the external application state is exercised
through a fake `App` and a `getenv` seam), gated in CI. Everything, including the
`rails/all` aggregate, cross-compiles on all six supported 64-bit targets:
`amd64`, `arm64`, `riscv64`, `loong64`, `ppc64le`, and `s390x`.

```console
$ go test -race -cover ./...
ok  github.com/go-ruby-rails/rails  coverage: 100.0% of statements
```

## License

BSD-3-Clause — see [LICENSE](LICENSE). Copyright (c) 2026, the
go-ruby-rails/rails authors.

## WebAssembly

Being pure Go (CGO=0), this library also compiles to **WebAssembly** — both
`GOOS=js GOARCH=wasm` (browser / Node.js) and `GOOS=wasip1 GOARCH=wasm` (WASI).
CI builds both targets on every push, alongside the six 64-bit native/qemu arches.

```sh
GOOS=js     GOARCH=wasm go build ./...   # browser / Node
GOOS=wasip1 GOARCH=wasm go build ./...   # WASI (wasmtime, wasmer, wasmedge, …)
```
