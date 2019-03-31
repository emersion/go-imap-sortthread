// Package sortthread implements SORT and THREAD for go-imap.
//
// SORT and THREAD are defined in RFC 5256.
package sortthread

const SortCapability = "SORT"

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
