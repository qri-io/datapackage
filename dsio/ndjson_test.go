package dsio

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/qri-io/dataset"
	"github.com/qri-io/dataset/compression"
)

func TestNDJSONReadWrite(t *testing.T) {
	data := `["a","b","c"]
"apples"
true
35
null
{}
`

	st := &dataset.Structure{
		Format: dataset.NDJSONDataFormat.String(),
		Schema: dataset.BaseSchemaArray,
	}

	rdr, err := NewEntryReader(st, strings.NewReader(data))
	if err != nil {
		t.Fatal(err)
	}

	buf := &bytes.Buffer{}
	wr, err := NewEntryWriter(st, buf)

	if err := Copy(rdr, wr); err != nil {
		t.Fatal(err)
	}
	rdr.Close()
	wr.Close()

	if diff := cmp.Diff(data, buf.String()); diff != "" {
		t.Errorf("result mismatch (-want +got):\n%s", diff)
	}
}

func TestNDJSONCompression(t *testing.T) {
	invalidCompressionSt := &dataset.Structure{Format: "ndjson", Compression: "invalid", Schema: dataset.BaseSchemaArray}
	if _, err := NewJSONReader(invalidCompressionSt, nil); err == nil {
		t.Errorf("constructing reader with invalid compression should error")
	}
	if _, err := NewJSONWriter(invalidCompressionSt, nil); err == nil {
		t.Errorf("constructing writer with invalid compression should error")
	}

	data := `["a","b","c"]
"apples"
true
35
null
{}
`

	compressed := &bytes.Buffer{}
	compressor, _ := compression.Compressor("zst", compressed)
	io.Copy(compressor, strings.NewReader(data))
	compressor.Close()

	st := &dataset.Structure{
		Format:      "ndjson",
		Compression: "zst",
		Schema:      dataset.BaseSchemaArray,
	}

	rdr, err := NewNDJSONReader(st, compressed)
	if err != nil {
		t.Fatal(err)
	}

	compressed2 := &bytes.Buffer{}
	wr, err := NewNDJSONWriter(st, compressed2)
	if err != nil {
		t.Fatal(err)
	}

	if err := Copy(rdr, wr); err != nil {
		t.Fatal(err)
	}
	rdr.Close()
	wr.Close()

	if diff := cmp.Diff(compressed.Bytes(), compressed2.Bytes()); diff != "" {
		t.Errorf("result mismatch expect (-want +got):\n%s", diff)
	}
}

func TestNDJSONReaderSizeOverflow(t *testing.T) {
	// run a test with one 24,000-character long string to ensure the reader
	// doesn't choke on a long line of JSON
	st := &dataset.Structure{
		Format: "ndjson",
		Schema: dataset.BaseSchemaArray,
	}
	data := fmt.Sprintf(`"hi"
false
%q
null
"bye"
`, strings.Repeat("long", 1024*6))

	rdr, err := NewNDJSONReader(st, strings.NewReader(data))
	if err != nil {
		t.Fatal(err)
	}

	vals, err := ReadAll(rdr)
	if err != nil {
		t.Error(err)
	}

	if err := rdr.Close(); err != nil {
		t.Error(err)
	}

	expect := []interface{}{
		"hi",
		false,
		strings.Repeat("long", 1024*6),
		nil,
		"bye",
	}

	if diff := cmp.Diff(expect, vals); diff != "" {
		t.Errorf("result mismatch (-want +got):\n%s", diff)
	}
}
