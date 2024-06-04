package quickwit

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

type Client struct {
	cfg     Config
	once    sync.Once
	quit    chan struct{}
	records chan record
	wg      sync.WaitGroup
}

type record map[string]any

// New makes a new Client from config
func New(cfg Config) *Client {
	if cfg.URL == "" {
		panic("missing Quickwit URL")
	}

	if cfg.BatchWait < 0 {
		panic("wrong BatchWait option")
	}
	if cfg.BatchBytes <= 0 {
		panic("wrong BatchWait option")
	}
	if cfg.Commit != Auto && cfg.Commit != WaitFor && cfg.Commit != Force {
		panic("wrong Commit option")
	}

	if cfg.Timeout <= 0 {
		panic("wrong Timeout option")
	}

	client := &Client{
		cfg:     cfg,
		once:    sync.Once{},
		quit:    make(chan struct{}),
		records: make(chan record),
		wg:      sync.WaitGroup{},
	}

	client.cfg.Client.Timeout = cfg.Timeout

	go client.run()

	return client
}

// NewWithDefault creates a new client with default configuration.
func NewWithDefault(url string) *Client {
	return New(NewDefaultConfig(url))
}

func (c *Client) run() {
	currentBatch := newBatch()
	c.wg.Add(1)

	minWaitCheckFrequency := 10 * time.Millisecond
	maxWaitCheckFrequency := c.cfg.BatchWait / 10
	if maxWaitCheckFrequency < minWaitCheckFrequency {
		maxWaitCheckFrequency = minWaitCheckFrequency
	}

	maxWaitCheck := time.NewTicker(maxWaitCheckFrequency)

	for {
		select {
		case <-c.quit:
			// send pending batch
			c.sendBatch(currentBatch)
			return

		case r := <-c.records:
			// check if the batch is full
			if currentBatch.size() > c.cfg.BatchBytes {
				currentBatch = c.rolloutBatches(currentBatch)
			}

			_ = currentBatch.add(r)

		case <-maxWaitCheck.C:
			// check if max wait time has been reached
			if currentBatch.age() < c.cfg.BatchWait {
				continue
			}

			currentBatch = c.rolloutBatches(currentBatch)
		}
	}
}

func (c *Client) rolloutBatches(previous *batch) *batch {
	c.wg.Add(1)
	go c.sendBatch(previous)
	return newBatch()
}

func (c *Client) sendBatch(batch *batch) {
	if batch.size() > 0 {
		backoff := newBackoff(context.Background(), c.cfg.BackoffConfig)

		for backoff.ongoing() {
			status, err := c.send(batch.buf)
			if err == nil {
				break
			}

			// Only retry 429s, 5xx and connection-level errors.
			if status > 0 && status != 429 && status < 500 {
				break
			}

			backoff.wait()
		}
	}

	c.wg.Done()
}

func (c *Client) send(buf *bytes.Buffer) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.cfg.Timeout)
	defer cancel()

	url := fmt.Sprintf("%s/api/v1/%s/ingest?commit=%s", string(c.cfg.Commit), c.cfg.IndexID, c.cfg.URL)

	req, err := http.NewRequestWithContext(ctx, "POST", url, buf)
	if err != nil {
		return 0, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "samber/go-qwikwit")

	resp, err := c.cfg.Client.Do(req)
	if err != nil {
		return 0, err
	}

	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return resp.StatusCode, fmt.Errorf("quickwit returned error %s (%d): %s", resp.Status, resp.StatusCode, body)
	}

	return resp.StatusCode, err
}

// Push adds a new record to the next batch. Delivery is async.
func (c *Client) Push(record record) {
	c.records <- record
}

// Stop the client.
func (c *Client) Stop() {
	c.once.Do(func() { close(c.quit) })
	c.wg.Wait()
}
