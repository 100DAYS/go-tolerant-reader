package tolerantreader

import (
	"encoding/json"
	"testing"
	"time"
)

type TestMsg struct {
	A float64   `jsonpath:"$.expensive"`
	B []float64 `jsonpath:"$.store.book[?(@.price < 10.0)].price"`
	C []string  `jsonpath:"$.store.book[:].category"`
	D time.Time `jsonpath:"$.time"`
	E time.Time `jsonpath:"$.date"`
	F int       `jsonpath:"$.int"`
}

func TestReader(t *testing.T) {
	var dat map[string]interface{}

	err := json.Unmarshal([]byte(data), &dat)
	if err != nil {
		t.Errorf("Error Unmarshalling to Map: %s", err)
	}

	o := TestMsg{}

	err = Unmarshal(dat, &o)
	if err != nil {
		t.Errorf("Error while Reading: %s", err)
	}

	if o.A != 10.0 {
		t.Errorf("Field expensive should be 10.00, found %f", o.A)
	}
	if len(o.B) != 2 {
		t.Errorf("expected 2 items with price < 10, fount  %d", len(o.B))
	}
	if o.C[0] != "reference" || o.C[1] != "fiction" || o.C[2] != "fiction" || o.C[3] != "fiction" {
		t.Errorf("categories incorrect: %v", o.C)
	}

	t1, err := time.Parse(time.RFC3339, "2019-12-31T10:10:22+02:00")
	if err != nil {
		panic(err)
	}
	if !o.D.Equal(t1) {
		t.Errorf("time incorrect: %v != %v", o.D, t1)
	}

	t2, err := time.Parse(time.RFC3339, "2019-10-30T00:00:00+00:00")
	if err != nil {
		panic(err)
	}
	if !o.E.Equal(t2) {
		t.Errorf("date incorrect: %v != %v", o.E, t2)
	}

	if o.F != 23 {
		t.Errorf("int incorrect: %v", o.F)
	}
}
