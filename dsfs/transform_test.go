package dsfs

import (
	"encoding/json"
	"testing"

	"github.com/ipfs/go-datastore"
	"github.com/qri-io/cafs"
	"github.com/qri-io/dataset"
	"github.com/qri-io/jsonschema"
)

func TestLoadTransform(t *testing.T) {
	// TODO - restore
	// store := cafs.NewMapstore()
	// q := &dataset.AbstractTransform{Statement: "select * from whatever booooooo go home"}
	// a, err := SaveAbstractTransform(store, q, true)
	// if err != nil {
	// 	t.Errorf(err.Error())
	// 	return
	// }

	// if _, err = LoadTransform(store, a); err != nil {
	// 	t.Errorf(err.Error())
	// }
	// TODO - other tests & stuff
}

func TestTransformLoadAbstract(t *testing.T) {
	// store := datastore.NewMapDatastore()
	// TODO - finish dis test
}

func TestSaveTransform(t *testing.T) {
	dsa := dataset.NewDatasetRef(datastore.NewKey("/path/to/dataset/a"))
	dsa.Assign(&dataset.Dataset{Meta: &dataset.Meta{Title: "now dataset isn't empty"}})

	store := cafs.NewMapstore()
	q := &dataset.Transform{
		Syntax: "sweet syntax",
		Structure: &dataset.Structure{
			Format: dataset.CSVDataFormat,
			Schema: jsonschema.Must(`true`),
		},
		Resources: map[string]*dataset.Dataset{
			"a": dsa,
		},
	}

	key, err := SaveTransform(store, q, true)
	if err != nil {
		t.Error(err.Error())
		return
	}

	hash := "/map/QmbaGqpbMmdJwRS4spApiwtE4JcwqHpR7YhvbpgUx1fqXF"
	if hash != key.String() {
		t.Errorf("key mismatch: %s != %s", hash, key.String())
		return
	}

	expectedEntries := 2
	if len(store.(cafs.MapStore)) != expectedEntries {
		t.Errorf("invalid number of entries added to store: %d != %d", expectedEntries, len(store.(cafs.MapStore)))
		return
	}

	f, err := store.Get(datastore.NewKey(hash))
	if err != nil {
		t.Errorf("error getting dataset file: %s", err.Error())
		return
	}

	res := &dataset.Transform{}
	if err := json.NewDecoder(f).Decode(res); err != nil {
		t.Errorf("error decoding transform json: %s", err.Error())
		return
	}

	if !res.Structure.IsEmpty() {
		t.Errorf("expected stored transform.Structure to be a reference")
	}
	for name, ref := range res.Resources {
		if !ref.IsEmpty() {
			t.Errorf("expected stored transform reference '%s' to be empty", name)
		}
	}
}

func TestSaveAbstractTransform(t *testing.T) {
	dsa := dataset.NewDatasetRef(datastore.NewKey("/path/to/dataset/a"))
	dsa.Assign(&dataset.Dataset{Meta: &dataset.Meta{Title: "now dataset isn't empty "}})
	dsa.Structure = &dataset.Structure{
		Format: dataset.CSVDataFormat,
		Schema: jsonschema.Must(`true`),
	}

	store := cafs.NewMapstore()
	q := &dataset.Transform{
		Syntax: "sweet syntax",
		Structure: &dataset.Structure{
			Format: dataset.CSVDataFormat,
			Schema: jsonschema.Must(`true`),
		},
		Resources: map[string]*dataset.Dataset{
			"a": dsa,
		},
	}

	key, err := SaveAbstractTransform(store, q, true)
	if err != nil {
		t.Error(err.Error())
		return
	}

	hash := "/map/QmXcMVfHAjhiahiXKEvAb9AvDahv4T7cG4FoVgkefScKEa"
	if hash != key.String() {
		t.Errorf("key mismatch: %s != %s", hash, key.String())
		return
	}

	expectedEntries := 3
	if len(store.(cafs.MapStore)) != expectedEntries {
		t.Errorf("invalid number of entries added to store: %d != %d", expectedEntries, len(store.(cafs.MapStore)))
		return
	}

	f, err := store.Get(datastore.NewKey(hash))
	if err != nil {
		t.Errorf("error getting dataset file: %s", err.Error())
		return
	}

	res := &dataset.Transform{}
	if err := json.NewDecoder(f).Decode(res); err != nil {
		t.Errorf("error decoding transform json: %s", err.Error())
		return
	}

	if !res.Structure.IsEmpty() {
		t.Errorf("expected stored transform.Structure to be a reference")
	}
	for name, ref := range res.Resources {
		if !ref.IsEmpty() {
			t.Errorf("expected stored transform reference '%s' to be empty", name)
		}
	}
}
