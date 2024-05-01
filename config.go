package quickwit

import (
	"net/http"
	"time"
)

type CommitMode string

const (
	Auto    CommitMode = "auto"
	WaitFor CommitMode = "wait_for"
	Force   CommitMode = "force"
)

const (
	DefaultBatchWait  time.Duration = 3 * time.Second
	DefaultBatchBytes int           = 10 * 1024 * 1024 // 10 MB
	DefautCommitMode                = Auto

	DefaultMinBackoff time.Duration = 500 * time.Millisecond
	DefaultMaxBackoff time.Duration = 5 * time.Minute
	DefaultMaxRetries int           = 5

	DefaultTimeout time.Duration = 10 * time.Second
)

// Config describes configuration for a HTTP pusher client.
type Config struct {
	URL    string
	Client http.Client

	BatchWait  time.Duration
	BatchBytes int
	Commit     CommitMode

	BackoffConfig BackoffConfig
	Timeout       time.Duration
}

// NewDefaultConfig creates a default configuration for a given Quickwit cluster.
func NewDefaultConfig(url string) Config {
	return Config{
		URL:    url,
		Client: http.Client{},

		BatchWait:  DefaultBatchWait,
		BatchBytes: DefaultBatchBytes,
		Commit:     DefautCommitMode,

		BackoffConfig: BackoffConfig{
			MinBackoff: DefaultMinBackoff,
			MaxBackoff: DefaultMaxBackoff,
			MaxRetries: DefaultMaxRetries,
		},
		Timeout: DefaultTimeout,
	}
}
