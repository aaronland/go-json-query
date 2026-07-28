package query

import (
	"context"
	"fmt"
	"github.com/tidwall/gjson"
	"log/slog"
	"regexp"
)

// QUERYSET_MODE_ANY is a flag to signal that only one match in a QuerySet needs to be successful.
const QUERYSET_MODE_ANY string = "ANY"

// QUERYSET_MODE_ALL is a flag to signal that only all matches in a QuerySet needs to be successful.
const QUERYSET_MODE_ALL string = "ALL"

// QuerySet is a struct containing one or more Query instances and flags for how the results of those queries should be interpreted.
type QuerySet struct {
	// A set of Query instances
	Queries []*Query
	// A string flag representing how query results should be interpreted.
	Mode string
}

func (qs *QuerySet) String() string {
	return fmt.Sprintf("%s %v", qs.Mode, qs.Queries)
}

// Query is an atomic query to perform against a JSON document.
type Query struct {
	// A valid tidwall/gjson query path.
	Path string
	// A valid regular expression.
	Match *regexp.Regexp
}

func (q *Query) String() string {
	return fmt.Sprintf("%s == %v", q.Path, q.Match)
}

// Matches compares the set of queries in 'qs' against a JSON record ('body') and returns true or false depending on whether or not some or all of those queries are matched successfully.
func Matches(ctx context.Context, qs *QuerySet, body []byte) (bool, error) {

	select {
	case <-ctx.Done():
		return false, nil
	default:
		// pass
	}

	queries := qs.Queries
	mode := qs.Mode

	tests := len(queries)
	matches := 0

	logger := slog.Default()
	logger = logger.With("query set", qs.String())
	
	for _, q := range queries {

		logger := slog.Default()
		// logger = logger.With("query set", qs.String())		
		logger = logger.With("path", q.Path, "re", q.Match)

		rsp := gjson.GetBytes(body, q.Path)

		if !rsp.Exists() {

			logger.Info("MISSING path")
			
			if mode == QUERYSET_MODE_ALL {
				logger.Info("Query mode is all, BREAK")
				break
			}
		}

		for _, r := range rsp.Array() {

			// logger.Info("TEST candidate", "index", idx, "candidate", r.String())
			
			if q.Match.MatchString(r.String()) {

				logger.Info("String MATCHES", "candidate", r.String())

				matches += 1

				if mode == QUERYSET_MODE_ANY {
					logger.Info("Query mode is any, BREAK")					
					break
				}
				
			} else {
				logger.Info("Match FAILED", "candidate", r.String())
			}
		}

		if mode == QUERYSET_MODE_ANY && matches > 0 {
			logger.Info("Matches is > 0 and query mode is any, BREAK")								
			break
		}
	}

	if mode == QUERYSET_MODE_ALL {

		if matches < tests {
			logger.Info("Return FALSE, query mode is all and matches < test", "tests", tests, "matches", matches)
			return false, nil
		}
	}

	if matches == 0 {
		logger.Info("Return FALSE, matches is 0")		
		return false, nil
	}

	logger.Info("Return TRUE")
	return true, nil
}
