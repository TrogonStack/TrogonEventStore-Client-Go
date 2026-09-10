# TrogonEventStore Go client

[![CI](https://github.com/TrogonStack/TrogonEventStore-Client-Go/actions/workflows/ci.yml/badge.svg)](https://github.com/TrogonStack/TrogonEventStore-Client-Go/actions/workflows/ci.yml)
[![Integration](https://github.com/TrogonStack/TrogonEventStore-Client-Go/actions/workflows/integration.yml/badge.svg)](https://github.com/TrogonStack/TrogonEventStore-Client-Go/actions/workflows/integration.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/TrogonStack/TrogonEventStore-Client-Go.svg)](https://pkg.go.dev/github.com/TrogonStack/TrogonEventStore-Client-Go/trogoneventstore)

The Go client for [TrogonEventStore](https://github.com/TrogonStack/TrogonEventStore).

This project is derived from an Apache-2.0 licensed client. The license retains
the original copyright notice and identifies modifications by Straw Hat, LLC.

## Install

```sh
go get github.com/TrogonStack/TrogonEventStore-Client-Go@latest
```

```go
import "github.com/TrogonStack/TrogonEventStore-Client-Go/trogoneventstore"
```

See [INSTALL.md](INSTALL.md) and the [API guides](docs/api/getting-started.md).

## Build and test

```sh
make build
go test ./trogoneventstore
```

Integration tests run against `ghcr.io/trogonstack/trogoneventstore:ci`:

```sh
make start-server
go test ./test
make stop-server
```

Override the server image with `TROGON_EVENTSTORE_SERVER_IMAGE` when testing a
specific image.

## Releases

Release Please creates version pull requests and GitHub releases. Go consumers
use the generated `vX.Y.Z` module tags directly. This project does not publish
to a package registry.

## License

Licensed under Apache-2.0. See [LICENSE](LICENSE).
