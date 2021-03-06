package dsio

import (
	"fmt"
	"io"

	"github.com/qri-io/dataset"
)

// PagedReader wraps a reader, starting reads from offset, and only reads limit number of entries
type PagedReader struct {
	Reader EntryReader
	Limit  int
	Offset int
}

var _ EntryReader = (*PagedReader)(nil)

// Structure returns the wrapped reader's structure
func (r *PagedReader) Structure() *dataset.Structure {
	return r.Reader.Structure()
}

// ReadEntry returns an entry, taking offset and limit into account
func (r *PagedReader) ReadEntry() (Entry, error) {
	for r.Offset > 0 {
		_, err := r.Reader.ReadEntry()
		if err != nil {
			return Entry{}, err
		}
		r.Offset--
	}
	if r.Limit == 0 {
		return Entry{}, io.EOF
	}
	r.Limit--
	return r.Reader.ReadEntry()
}

// Close finalizes the writer, indicating no more records
// will be written
func (r *PagedReader) Close() error {
	return r.Reader.Close()
}

// Copy reads all entries from the reader and writes them to the writer
func Copy(reader EntryReader, writer EntryWriter) error {
	for {
		val, err := reader.ReadEntry()
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("row iteration error: %s", err.Error())
		}
		if err := writer.WriteEntry(val); err != nil {
			return fmt.Errorf("error writing value to buffer: %s", err.Error())
		}
	}
	return nil
}
