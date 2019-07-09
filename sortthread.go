// Package sortthread implements SORT and THREAD for go-imap.
//
// SORT and THREAD are defined in RFC 5256.
package sortthread

import (
	"fmt"
	"strconv"
)

const SortCapability = "SORT"

var ThreadCapabilities = []string{"THREAD=ORDEREDSUBJECT", "THREAD=REF", "THREAD=REFERENCES"}

// ThreadAlgorithm is the algorithm used by the server to sort messages
type ThreadAlgorithm string

const (
	OrderedSubject ThreadAlgorithm = "ORDEREDSUBJECT"
	References                     = "REFERENCES"
)

// SortField is a field that can be used to sort messages.
type SortField string

const (
	SortArrival SortField = "ARRIVAL"
	SortCc = "CC"
	SortDate = "DATE"
	SortFrom = "FROM"
	SortSize = "SIZE"
	SortSubject = "SUBJECT"
	SortTo = "TO"
)

// SortCriterion is a criterion that can be used to sort messages.
type SortCriterion struct {
	Field SortField
	Reverse bool
}

func formatSortCriteria(criteria []SortCriterion) interface{} {
	fields := make([]interface{}, 0, len(criteria))
	for _, c := range criteria {
		if c.Reverse {
			fields = append(fields, "REVERSE")
		}
		fields = append(fields, string(c.Field))
	}
	return fields
}

type Thread struct {
	Id       uint32
	Children []*Thread
}

func (t *Thread) String() string {
	return fmt.Sprintf("(%s)", t.toString())
}

func (t *Thread) toString() string {
	s := strconv.FormatUint(uint64(t.Id), 10)
	if len(t.Children) == 1 {
		s += fmt.Sprintf(" %s", t.Children[0].toString())
	} else {
		for _, child := range t.Children {
			s += fmt.Sprintf(" (%s)", child.toString())
		}
	}
	return s
}
