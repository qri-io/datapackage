package dstest

import (
	"bytes"
	"testing"
)

func TestNewTestCaseFromDir(t *testing.T) {
	tc, err := NewTestCaseFromDir("testdata/complete", t)
	if err != nil {
		t.Errorf("error reading test dir: %s", err.Error())
		return
	}

	name := "complete"
	if tc.Name != name {
		t.Errorf("expected name to equal: %s. got: %s", name, tc.Name)
	}

	fn := "data.csv"
	if tc.DataFilename != fn {
		t.Errorf("expected DataFilename to equal: %s. got: %s", fn, tc.DataFilename)
	}

	data := []byte(`city,pop,avg_age,in_usa
toronto,40000000,55.5,false
new york,8500000,44.4,true
chicago,300000,44.4,true
chatham,35000,65.25,true
raleigh,250000,50.65,true
`)
	if !bytes.Equal(tc.Data, data) {
		t.Errorf("data mismatch")
	}
}