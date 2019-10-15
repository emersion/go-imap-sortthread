// Package sortthread implements SORT and THREAD for go-imap.
//
// SORT and THREAD are defined in RFC 5256.
package sortthread

import (
	"fmt"
	"mime"
	"regexp"
	"strings"

	"github.com/emersion/go-imap"
)

const SortCapability = "SORT"

var ThreadCapabilities = []string{"THREAD=ORDEREDSUBJECT", "THREAD=REF", "THREAD=REFERENCES"}

// ThreadAlgorithm is the algorithm used by the server to sort messages
type ThreadAlgorithm string

const (
	OrderedSubject ThreadAlgorithm = "ORDEREDSUBJECT"
	References                     = "REFERENCES"
)

func formatThreadAlgorithm(algorithm ThreadAlgorithm) imap.RawString {
	return imap.RawString(algorithm)
}

// SortField is a field that can be used to sort messages.
type SortField string

const (
	SortArrival SortField = "ARRIVAL"
	SortCc                = "CC"
	SortDate              = "DATE"
	SortFrom              = "FROM"
	SortSize              = "SIZE"
	SortSubject           = "SUBJECT"
	SortTo                = "TO"
)

// SortCriterion is a criterion that can be used to sort messages.
type SortCriterion struct {
	Field   SortField
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

var (
	tabsContinuation = regexp.MustCompile(`[\t\\]`)
	repeatedSpaces   = regexp.MustCompile("[ ]+")

	// Includes regex for ABNF rules relevant to base subject
	// Note that all ABNF strings are considered lowercase
	subjFwdHeader = "[fwd:"
	subjFwdTrail  = "]"

	// BLOBCHAR        = %x01-5a / %x5c / %x5e-ff
	// subj-blob       = "[" *BLOBCHAR "]" *WSP
	subjBlob       = `\[\x01-\x5a\x5c\x5e-\xff]\]\s*`
	subjBlobPrefix = regexp.MustCompile(fmt.Sprintf("^%s", subjBlob))

	// subj-refwd      = ("re" / ("fw" ["d"])) *WSP [subj-blob] ":"
	subjReFwd = fmt.Sprintf(`(?:(?:re)|(?:fwd?))\s*(?:%s)?:`, subjBlob)
	// subj-leader     = (*subj-blob subj-refwd) / WSP
	subjLeader = regexp.MustCompile(fmt.Sprintf(`(?i)^(?:(?:%s)*%s)`,
		subjBlob, subjReFwd))
	// subj-trailer    = "(fwd)" / WSP
	subjTrailer = regexp.MustCompile(`(?i)\(fwd\)$`)
)

// Steps 2-4 in RFC Section 2.1
func replaceArtifacts(subject string, isReplyFwd *bool) string {
	// (2) Remove all trailing text of the subject that matches the
	// subj-trailer ABNF; repeat until no more matches are possible.
	for {
		noTrail := strings.TrimSuffix(subject, " ")
		if subjTrailer.MatchString(noTrail) {

			noTrail = subjTrailer.ReplaceAllString(noTrail, "")
			*isReplyFwd = true
		}
		if subject == noTrail {
			break
		}
		subject = noTrail
	}
	return replacePrefix(subject, isReplyFwd)
}

// Steps 3 & 4 in RFC Section 2.1
func replacePrefix(subject string, isReplyFwd *bool) string {
	// Repeat (3) and (4) until no matches remain.
	for {
		// (3) Remove all prefix text of the subject that matches the subj-
		// leader ABNF.
		noLeader := strings.TrimPrefix(subject, " ")
		if subjLeader.MatchString(noLeader) {
			noLeader = subjLeader.ReplaceAllString(noLeader, "")
			*isReplyFwd = true
		}

		// (4) If there is prefix text of the subject that matches the subj-
		// blob ABNF, and removing that prefix leaves a non-empty subj-
		// base, then remove the prefix text.
		noBlob := subjBlobPrefix.ReplaceAllString(noLeader, "")
		if noBlob == "" {
			subject = noLeader
			break
		}
		if noBlob == subject {
			break
		}
		subject = noBlob
	}
	return subject
}

func GetBaseSubject(subject string, isReplyFwd *bool) (string, error) {
	var (
		baseSubject string
		err         error
	)

	*isReplyFwd = false

	// (1) Convert any RFC 2047 encoded-words in the subject to [UTF-8]
	// as described in "Internationalization Considerations".
	// Convert all tabs and continuations to space.  Convert all
	// multiple spaces to a single space.
	dec := new(mime.WordDecoder)
	baseSubject, err = dec.DecodeHeader(subject)
	if err != nil {
		return "", err
	}

	baseSubject = tabsContinuation.ReplaceAllString(baseSubject, " ")
	baseSubject = repeatedSpaces.ReplaceAllString(baseSubject, " ")

	for {
		// Steps 2-5
		baseSubject = replaceArtifacts(baseSubject, isReplyFwd)

		// (6) If the resulting text begins with the subj-fwd-hdr ABNF and
		// ends with the subj-fwd-trl ABNF, remove the subj-fwd-hdr and
		// subj-fwd-trl and repeat from step (2).
		if !strings.HasPrefix(baseSubject, subjFwdHeader) || !strings.HasSuffix(baseSubject, subjFwdTrail) {
			break
		}
		*isReplyFwd = true
		baseSubject = strings.TrimPrefix(baseSubject, subjFwdHeader)
		baseSubject = strings.TrimSuffix(baseSubject, subjFwdTrail)
	}

	return baseSubject, nil
}
