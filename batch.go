package quickwit

import (
	"bytes"
	"encoding/json"
	"time"
)

// batch holds pending items waiting to be sent to Quickwit, and it's used
// to reduce the number of push requests to Quickwit aggregating multiple records
// in a single batch request.
type batch struct {
	buf       *bytes.Buffer
	createdAt time.Time
}

func newBatch(records ...record) *batch {
	b := &batch{
		buf:       bytes.NewBuffer([]byte{}),
		createdAt: time.Now(),
	}

	// Add entries to the batch
	for _, record := range records {
		_ = b.add(record)
	}

	return b
}

// add an entry to the batch
func (b *batch) add(record record) error {
	bytes, err := json.Marshal(record)
	if err != nil {
		return err
	}

	b.buf.Write(bytes)
	b.buf.WriteByte('\n')
	return nil
}

// size returns the current batch size in bytes
func (b *batch) size() int {
	return b.buf.Len()
}

// age of the batch since its creation
func (b *batch) age() time.Duration {
	return time.Since(b.createdAt)
}
