package query

import (
	"testing"
	"context"
)

func TestMatches(t *testing.T) {

	ctx := context.Background()
	
	doc := []byte(`{"properties": { "wof:placetype": "campus" }`)
	
	flags := []string{
		"properties.wof:placetype=campus",
		"properties.wof:placetype=locality",
	}
	
	qf := new(QueryFlags)

	for _, fl := range flags {
		
		err := qf.Set(fl)

		if err != nil {
			t.Fatalf("Failed to set flag '%s', %v", fl, err)
		}
	}

	qs := &QuerySet{
		Queries: *qf,
		Mode: "ANY",
	}

	matches, err := Matches(ctx, qs, doc)

	if err != nil {
		t.Fatalf("Failed to match, %v", err)
	}

	if !matches {
		t.Fatalf("Invalid match in ANY mode")
	}

	qs2 := &QuerySet{
		Queries: *qf,
		Mode: "ALL",
	}

	matches2, err := Matches(ctx, qs2, doc)

	if err != nil {
		t.Fatalf("Failed to match, %v", err)
	}

	if matches2 {
		t.Fatalf("Invalid match in ALL mode")
	}
	
}
