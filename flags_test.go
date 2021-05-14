package query

import (
	"testing"
)

func TestQueryFlags(t *testing.T) {

	ok_flags := []string{
		"properties.wof:placetype=campus",
		"properties.mz_is_current=^(1|0)$",
	}
	
	qf := new(QueryFlags)

	for _, f := range ok_flags {
		
		err := qf.Set(f)

		if err != nil {
			t.Fatalf("Failed to set flag '%s', %v", f, err)
		}
	}
}
