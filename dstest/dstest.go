// Package dstest defines an interface for reading test cases from static files
// leveraging directories of test dataset input files & expected output files
package dstest

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/qri-io/cafs"
	"github.com/qri-io/dataset"
)

const (
	// InputDatasetFilename is the filename to use for an input dataset
	InputDatasetFilename = "input.dataset.json"
	// ExpectDatasetFilename is the filename to use to compare expected outputs
	ExpectDatasetFilename = "expect.dataset.json"
)

// TestCase is a dataset test case, usually built from a
// directory of files for use in tests.
// All files are optional for TestCase, but may be required
// by the test itself.
type TestCase struct {
	// Name is the casename, should match directory name
	Name string
	// 	 data.csv,data.json, etc
	DataFilename string
	// test data in expected data format
	Data []byte
	// Input is intended file for test input
	// loads from input.dataset.json
	Input *dataset.Dataset
	//  Expect should match test output
	// loads from expect.dataset.json
	Expect *dataset.Dataset
}

// DataFile creates a new in-memory file from data & filename properties
func (t TestCase) DataFile() cafs.File {
	return cafs.NewMemfileBytes(t.DataFilename, t.Data)
}

// NewTestCaseFromDir creates a test case from a directory of static test files
// dir should be the path to the directory to check, and any parsing errors will
// be logged using t.Log methods
func NewTestCaseFromDir(dir string, t *testing.T) (tc TestCase, err error) {
	tc = TestCase{
		Name: filepath.Base(dir),
	}
	tc.Data, tc.DataFilename, err = ReadInputData(dir)
	if err != nil {
		err = fmt.Errorf("error reading test case data for directory %s: %s", dir, err.Error())
		return
	}

	tc.Input, err = ReadDataset(dir, InputDatasetFilename)
	if err != nil && !os.IsNotExist(err) && t != nil {
		t.Logf("%s: error loading input dataset: %s", tc.Name, err)
	}
	err = nil

	tc.Expect, err = ReadDataset(dir, ExpectDatasetFilename)
	if err != nil && !os.IsNotExist(err) && t != nil {
		t.Logf("%s: error loading expect dataset: %s", tc.Name, err)
	}
	err = nil

	return
}

// ReadDataset grabs a dataset for a given dir for a given filename
func ReadDataset(dir, filename string) (*dataset.Dataset, error) {
	data, err := ioutil.ReadFile(filepath.Join(dir, filename))
	if err != nil {
		return nil, err
	}

	ds := &dataset.Dataset{}
	return ds, ds.UnmarshalJSON(data)
}

// ReadInputData grabs input data
func ReadInputData(dir string) ([]byte, string, error) {
	for _, df := range dataset.SupportedDataFormats() {
		path := fmt.Sprintf("%s/data.%s", dir, df)
		if f, err := os.Open(path); err == nil {
			data, err := ioutil.ReadAll(f)
			return data, fmt.Sprintf("data.%s", df), err
		}
	}
	return nil, "", os.ErrNotExist
}
