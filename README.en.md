# Kite

`kite` - Personal developer tool command application.

![app cmds](docs/images/kite-in-wsl.png)

## Features

* Git common command operations
* GitLab common command operation
* GitHub common command operations
* String processing tools: analysis, formatting, extracting information, converting
* json processing tools: formatting, searching, filtering, etc.
* go, php, java code generation, conversion, etc.
* json, yaml, sql formatting, conversion
* View system and environment information
* Quickly run built-in scripts
* Quickly run system commands
* File search matching, processing
* Batch run commands
* Document search, viewing, etc.

## Git Repos

* [https://github.com/inhere/kite-go](https://github.com/inhere/kite-go) Go language version
* [https://github.com/inhere/kite](https://github.com/inhere/kite)  PHP version

## Install

### Quick install

```bash
curl https://raw.githubusercontent.com/inhere/kite-go/main/cmd/install.sh | bash
```

**From proxy**

```shell
curl https://ghproxy.com/https://raw.githubusercontent.com/inhere/kite-go/main/cmd/install.sh | bash
```

### Install by go

```bash
go install github.com/inhere/kite-go/cmd/kite
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
- https://github.com/inhere/kite
