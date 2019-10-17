package sortthread

import (
	"errors"

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

func (cmd *SortCommand) Parse(fields []interface{}) error {
	return errors.New("sortthread: not yet implemented")
}

// ThreadCommand is a THREAD command.
type ThreadCommand struct {
	Algorithm      ThreadAlgorithm
	Charset        string
	SearchCriteria *imap.SearchCriteria
}

func (cmd *ThreadCommand) Command() *imap.Command {
	args := []interface{}{
		formatThreadAlgorithm(cmd.Algorithm),
		cmd.Charset,
	}

	addAllSearchCriteria := true

	// Verify if SearchCriteria is empty to use "ALL" as criteria
	if cmd.SearchCriteria != nil {
		searchFormat := cmd.SearchCriteria.Format()

		if len(searchFormat) > 0 {
			addAllSearchCriteria = false

			args = append(args, searchFormat)
		}
	}

	if addAllSearchCriteria {
		args = append(args, imap.RawString("ALL"))
	}

	return &imap.Command{
		Name:      "THREAD",
		Arguments: args,
	}
}
