# Kite

`Kite` - A CLI tools package.

> Github: https://github.com/inhere/kite-go

## Features

- quick create new project
- hot reload serve on file modified
- generate simple/controller/restful codes
- install development tools. eg: swaggo, swaggerui, golint, revive

## PHP version

- [inhere/kite](https://github.com/inhere/kite)

## Install

### Quick install

```bash
curl https://raw.githubusercontent.com/inhere/kite-go/main/cmd/install.sh | bash
```

### Install by go

```bash
go install github.com/inhere/kite/cmd/kite
```

## Build

```bash
make install
# or
go build -o $GOPAHT/bin/kite ./cmd/kite
```

## Develop

### Dev build

```shell
KITE_INIT_LOG=debug go run ./cmd/kite
```

### Install to GOBIN

```bash
make kit2gobin
# or
make kite2gobin
```

## Gookit Packages

- https://github.com/gookit/config
- https://github.com/gookit/rux
- https://github.com/gookit/gcli
- https://github.com/gookit/ini

## Refers

- https://github.com/bitfield/script
