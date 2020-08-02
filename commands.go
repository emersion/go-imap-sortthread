package sortthread

import (
	"errors"
	"io"
	"strings"

	"github.com/emersion/go-imap"
)

// SortCommand is a SORT command.
type SortCommand struct {
	SortCriteria   []SortCriterion
	Charset        string
	SearchCriteria *imap.SearchCriteria
}

func (cmd *SortCommand) Command() *imap.Command {
	args := []interface{}{
		formatSortCriteria(cmd.SortCriteria),
		cmd.Charset,
	}
	args = append(args, cmd.SearchCriteria.Format()...)

	return &imap.Command{
		Name:      "SORT",
		Arguments: args,
	}
}

func parseSortCriteria(fields interface{}) ([]SortCriterion, error) {
	list, ok := fields.([]interface{})
	if !ok {
		return nil, errors.New("List is required as a sort criteria")
	}

	result := make([]SortCriterion, 0, len(list))
	reverse := false
	for _, crit := range list {
		crit, ok := crit.(string)
		if !ok {
			return nil, errors.New("String is required as a sort key")
		}

		if strings.EqualFold(crit, "REVERSE") {
			reverse = true
			continue
		}
		crit = strings.ToUpper(crit)
		switch crit {
		// TODO: Fix types for constants.
		case string(SortArrival), SortCc, SortDate, SortFrom, SortSize, SortSubject, SortTo:
		default:
			return nil, errors.New("Unknown sort criteria: " + crit)
		}
		result = append(result, SortCriterion{
			Field:   SortField(crit),
			Reverse: reverse,
		})
		reverse = false
	}

	if reverse {
		return nil, errors.New("Missing sort key after REVERSE")
	}

	return result, nil
}

func (cmd *SortCommand) Parse(fields []interface{}) error {
	if len(fields) < 3 {
		return errors.New("Not enough SORT arguments")
	}

	var err error
	cmd.SortCriteria, err = parseSortCriteria(fields[0])
	if err != nil {
		return err
	}

	// Charset parameter for SORT is specified without "CHARSET"
	// and is required.
	charset, ok := fields[1].(string)
	if !ok {
		return errors.New("String is required as a charset")
	}
	charset = strings.ToLower(charset)
	var charsetReader func(io.Reader) io.Reader
	if charset != "utf-8" && charset != "us-ascii" && charset != "" {
		charsetReader = func(r io.Reader) io.Reader {
			r, _ = imap.CharsetReader(charset, r)
			return r
		}
	}

	cmd.SearchCriteria = &imap.SearchCriteria{}
	return cmd.SearchCriteria.ParseWithCharset(fields[2:], charsetReader)
}

// ThreadCommand is a THREAD command.
type ThreadCommand struct {
	Algorithm      ThreadAlgorithm
	Charset        string
	SearchCriteria *imap.SearchCriteria
}

func (cmd *ThreadCommand) Command() *imap.Command {
	return &imap.Command{
		Name: "THREAD",
		Arguments: []interface{}{
			formatThreadAlgorithm(cmd.Algorithm),
			cmd.Charset,
			cmd.SearchCriteria.Format(),
		},
	}
}

func (cmd *ThreadCommand) Parse(fields []interface{}) error {
	if len(fields) < 3 {
		return errors.New("Not enough THREAD argments")
	}

	algo, ok := fields[0].(string)
	if !ok {
		return errors.New("First argument should be a string")
	}

	charset, ok := fields[1].(string)
	if !ok {
		return errors.New("Second argument should be a string")
	}
	charset = strings.ToLower(charset)
	var charsetReader func(io.Reader) io.Reader
	if charset != "utf-8" && charset != "us-ascii" && charset != "" {
		charsetReader = func(r io.Reader) io.Reader {
			r, _ = imap.CharsetReader(charset, r)
			return r
		}
	}

	cmd.Algorithm = ThreadAlgorithm(algo)
	cmd.SearchCriteria = &imap.SearchCriteria{}
	return cmd.SearchCriteria.ParseWithCharset(fields[2:], charsetReader)
}
