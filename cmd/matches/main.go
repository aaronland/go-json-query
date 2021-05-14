// Read one or more files and test whether their contents match one or more query parameters.
package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/aaronland/go-json-query"
	"io"
	"log"
	"os"
	"strings"
)

func main() {

	var queries query.QueryFlags
	flag.Var(&queries, "query", "One or more {PATH}={REGEXP} parameters for filtering records.")

	valid_modes := strings.Join([]string{query.QUERYSET_MODE_ALL, query.QUERYSET_MODE_ANY}, ", ")
	desc_modes := fmt.Sprintf("Specify how query filtering should be evaluated. Valid modes are: %s", valid_modes)

	query_mode := flag.String("query-mode", query.QUERYSET_MODE_ALL, desc_modes)

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage:\n")
		fmt.Fprintf(os.Stderr, "  %s [options] [path1 path2 ... pathN]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	paths := flag.Args()

	if len(queries) == 0 {
		log.Fatalf("Nothing to query!")
	}

	qs := &query.QuerySet{
		Queries: queries,
		Mode:    *query_mode,
	}

	ctx := context.Background()

	for _, path := range paths {

		fh, err := os.Open(path)

		if err != nil {
			log.Fatalf("Failed to open '%s', %v", path, err)
		}

		defer fh.Close()

		body, err := io.ReadAll(fh)

		if err != nil {
			log.Fatalf("Failed to read '%s', %v", path, err)
		}

		matches, err := query.Matches(ctx, qs, body)

		if err != nil {
			log.Fatalf("Failed to query '%s', %v", path, err)
		}

		fmt.Printf("%s\t%t\n", path, matches)
	}
}
