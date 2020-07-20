package sortthread

import (
	"strconv"

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
	fields := make([]interface{}, 0, len(r.Ids)+1)
	fields = append(fields, imap.RawString("SORT"))
	for _, id := range r.Ids {
		fields = append(fields, imap.RawString(strconv.FormatInt(int64(id), 10)))
	}

	return imap.NewUntaggedResp(fields).WriteTo(w)
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
		case string, imap.RawString:
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

func formatThread(thread *Thread) []interface{} {
	f := make([]interface{}, 0, 1+len(thread.Children))
	f = append(f, imap.RawString(strconv.FormatInt(int64(thread.Id), 10)))
	if len(thread.Children) == 1 {
		f = append(f, formatThread(thread.Children[0])...)
	} else {
		for _, c := range thread.Children {
			f = append(f, formatThread(c))
		}
	}
	return f
}

func formatThreadResp(threads []*Thread) []interface{} {
	fields := make([]interface{}, 0, len(threads)+1)
	fields = append(fields, imap.RawString("THREAD"))
	for _, t := range threads {
		fields = append(fields, formatThread(t))
	}
	return fields
}

func (r *ThreadResponse) WriteTo(w *imap.Writer) error {
	return imap.NewUntaggedResp(formatThreadResp(r.Threads)).WriteTo(w)
}
