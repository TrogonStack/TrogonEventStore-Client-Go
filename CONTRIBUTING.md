# Contributing to the TrogonEventStore Go client

Open an issue or pull request in this repository. Pull request titles and commit
subjects must use `feat`, `fix`, or `chore` Conventional Commit prefixes. Every
commit must include a DCO sign-off.

Before opening a pull request, run:

```sh
gofmt -w .
go vet ./...
go test ./...
```

Integration changes must also be tested against the repository Docker Compose
stack.
