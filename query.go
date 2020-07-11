package query

import (
	"github.com/tidwall/gjson"
	"regexp"
)

const QUERYSET_MODE_ANY string = "ANY"
const QUERYSET_MODE_ALL string = "ALL"

type QuerySet struct {
	Queries []*Query
	Mode    string
}

type Query struct {
	Path  string
	Match *regexp.Regexp
}

func Matches(body []byte, qs QuerySet) (bool, error) {

	queries := qs.Queries
	mode := qs.Mode

	tests := len(queries)
	matches := 0

	for _, q := range queries {

		rsp := gjson.GetBytes(body, q.Path)

		if !rsp.Exists() {

			if mode == QUERYSET_MODE_ALL {
				break
			}
		}

		for _, r := range rsp.Array() {

			has_match := true

			if !q.Match.MatchString(r.String()) {

				has_match = false

				if mode == QUERYSET_MODE_ALL {
					break
				}
			}

			if !has_match {

				if mode == QUERYSET_MODE_ALL {
					break
				}

				continue
			}

			matches += 1
		}
	}

	if mode == QUERYSET_MODE_ALL {

		if matches < tests {
			return false, nil
		}
	}

	if matches == 0 {
		return false, nil
	}

	return true, nil
}
