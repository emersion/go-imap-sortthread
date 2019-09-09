package sortthread

import (
	"errors"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/responses"
)

type SortResponse struct {
	Ids []uint32
}

type ThreadResponse struct {
	Threads []*Thread
}

func (r *SortResponse) Handle(resp imap.Resp) error {
	name, fields, ok := imap.ParseNamedResp(resp)
	if !ok || name != "SORT" {
		return responses.ErrUnhandled
	}

	r.Ids = make([]uint32, len(fields))
	for i, f := range fields {
		if id, err := imap.ParseNumber(f); err != nil {
			return err
		} else {
			r.Ids[i] = id
		}
	}

	return nil
}

func (r *SortResponse) WriteTo(w *imap.Writer) error {
	return errors.New("sortthread: not yet implemented")
}

func (r *ThreadResponse) Handle(resp imap.Resp) error {
	name, fields, ok := imap.ParseNamedResp(resp)
	if !ok || name != "THREAD" {
		return responses.ErrUnhandled
	}
	if threads, err := parseThreadResp(fields); err != nil {
		return err
	} else {
		r.Threads = threads
	}

	return nil
}

func parseThreadResp(fields []interface{}) ([]*Thread, error) {
	var parent *Thread
	var siblings []*Thread
	for _, f := range fields {
		switch f := f.(type) {
		case string:
			id, err := imap.ParseNumber(f)
			if err != nil {
				return nil, err
			}
			t := Thread{Id: id}
			if parent == nil {
				siblings = append(siblings, &t)
			} else {
				parent.Children = append(t.Children, &t)
			}
			parent = &t
		case []interface{}:
			t, err := parseThreadResp(f)
			if err != nil {
				return nil, err
			}
			// Parent doesn't exist, e.g. didn't match the search
			// criteria. Let's ignore the parent thread.
			if parent == nil {
				siblings = append(siblings, t...)
			} else {
				parent.Children = append(parent.Children, t...)
			}
		default:
			return nil, responses.ErrUnhandled
		}
	}
	return siblings, nil
}
