
# A Go ingestion client for Quickwit

[![tag](https://img.shields.io/github/tag/samber/go-quickwit.svg)](https://github.com/samber/go-quickwit/releases)
![Go Version](https://img.shields.io/badge/Go-%3E%3D%201.18.0-%23007d9c)
[![GoDoc](https://godoc.org/github.com/samber/go-quickwit?status.svg)](https://pkg.go.dev/github.com/samber/go-quickwit)
![Build Status](https://github.com/samber/go-quickwit/actions/workflows/test.yml/badge.svg)
[![Go report](https://goreportcard.com/badge/github.com/samber/go-quickwit)](https://goreportcard.com/report/github.com/samber/go-quickwit)
[![Coverage](https://img.shields.io/codecov/c/github/samber/go-quickwit)](https://codecov.io/gh/samber/go-quickwit)
[![Contributors](https://img.shields.io/github/contributors/samber/go-quickwit)](https://github.com/samber/go-quickwit/graphs/contributors)
[![License](https://img.shields.io/github/license/samber/go-quickwit)](./LICENSE)

A [Quickwit](https://quickwit.io/) push client for Go. See [slog-quickwit](https://github.com/samber/slog-quickwit/) for a slog handler implementation.

If you're looking for a search library or Quickwit management interface, check the [official library](https://github.com/quickwit-oss/quickwit-go).

## 🚀 Install

```sh
go get github.com/samber/go-quickwit
```

This library is v0 and follows SemVer strictly. Some breaking changes might be made to exported APIs before v1.0.0.

## 💡 Spec

GoDoc: [https://pkg.go.dev/github.com/samber/go-quickwit](https://pkg.go.dev/github.com/samber/go-quickwit)

```go
type Config struct {
	URL    string
	Client http.Client

	BatchWait  time.Duration
	BatchBytes int
	Commit     CommitMode   // either quickwit.Auto, quickwit.WaitFor or quickwit.Force

	BackoffConfig BackoffConfig
	Timeout       time.Duration
}

type BackoffConfig struct {
	// start backoff at this level
	MinBackoff time.Duration
	// increase exponentially to this level
	MaxBackoff time.Duration
	// give up after this many; zero means infinite retries
	MaxRetries int
}
```

## Example

First, start Quickwit:

```bash
docker-compose up -d
curl -X POST \
    'http://localhost:7280/api/v1/indexes' \
    -H 'Content-Type: application/yaml' \
    --data-binary @test-config.yaml
```

Then push logs:

```go
import "github.com/samber/go-quickwit"

func main() {
	client := quickwit.NewWithDefault("http://localhost:7280")
	defer client.Stop() // flush and stop

	for i := 0; i < 10; i++ {
		client.Push(map[string]any{
			"timestamp": time.Now().Unix(),
			"message":   fmt.Sprintf("hello %d", i),
		})
		time.Sleep(1 * time.Second)
	}
}
```

## 🤝 Contributing

- Ping me on Twitter [@samuelberthe](https://twitter.com/samuelberthe) (DMs, mentions, whatever :))
- Fork the [project](https://github.com/samber/go-quickwit)
- Fix [open issues](https://github.com/samber/go-quickwit/issues) or request new features

Don't hesitate ;)

```bash
# start quickwit
docker-compose up -d
curl -X POST \
    'http://localhost:7280/api/v1/indexes' \
    -H 'Content-Type: application/yaml' \
    --data-binary @test-config.yaml

# Install some dev dependencies
make tools

# Run tests
make test
# or
make watch-test
```

## 👤 Contributors

![Contributors](https://contrib.rocks/image?repo=samber/go-quickwit)

## 💫 Show your support

Give a ⭐️ if this project helped you!

[![GitHub Sponsors](https://img.shields.io/github/sponsors/samber?style=for-the-badge)](https://github.com/sponsors/samber)

## 📝 License

Copyright © 2024 [Samuel Berthe](https://github.com/samber).

This project is [MIT](./LICENSE) licensed.
