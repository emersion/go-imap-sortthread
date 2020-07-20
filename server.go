package sortthread

import (
	"errors"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/backend"
	"github.com/emersion/go-imap/server"
)

var ErrUnsupportedBackend = errors.New("sortthread: backend not supported")

type SortMailbox interface {
	backend.Mailbox
	Sort(uid bool, sortCrit []SortCriterion, searchCrit *imap.SearchCriteria) ([]uint32, error)
}

type ThreadBackend interface {
	backend.Backend
	SupportedThreadAlgorithms() []ThreadAlgorithm
}

type ThreadMailbox interface {
	backend.Mailbox
	Thread(uid bool, threading ThreadAlgorithm, searchCrit *imap.SearchCriteria) ([]*Thread, error)
}

type SortHandler struct {
	SortCommand
}

func (h *SortHandler) handle(uid bool, conn server.Conn) error {
	if conn.Context().Mailbox == nil {
		return server.ErrNoMailboxSelected
	}

	mbox, ok := conn.Context().Mailbox.(SortMailbox)
	if !ok {
		return ErrUnsupportedBackend
	}

	ids, err := mbox.Sort(uid, h.SortCriteria, h.SearchCriteria)
	if err != nil {
		return err
	}

	return conn.WriteResp(&SortResponse{Ids: ids})
}

func (h *SortHandler) Handle(conn server.Conn) error {
	return h.handle(false, conn)
}

func (h *SortHandler) UidHandle(conn server.Conn) error {
	return h.handle(true, conn)
}

type ThreadHandler struct {
	ThreadCommand
}

func (h *ThreadHandler) handle(uid bool, conn server.Conn) error {
	if conn.Context().Mailbox == nil {
		return server.ErrNoMailboxSelected
	}

	mbox, ok := conn.Context().Mailbox.(ThreadMailbox)
	if !ok {
		return ErrUnsupportedBackend
	}

	thr, err := mbox.Thread(uid, h.Algorithm, h.SearchCriteria)
	if err != nil {
		return err
	}

	return conn.WriteResp(&ThreadResponse{Threads: thr})
}

func (h *ThreadHandler) Handle(conn server.Conn) error {
	return h.handle(false, conn)
}

func (h *ThreadHandler) UidHandle(conn server.Conn) error {
	return h.handle(true, conn)
}

type sortExtension struct{}

func NewSortExtension() server.Extension {
	return &sortExtension{}
}

func (s *sortExtension) Capabilities(c server.Conn) []string {
	if c.Context().State&imap.AuthenticatedState != 0 {
		return []string{SortCapability}
	}
	return nil
}

func (s *sortExtension) Command(name string) server.HandlerFactory {
	if name == "SORT" {
		return func() server.Handler {
			return &SortHandler{}
		}
	}
	return nil
}

type threadExtension struct{}

func NewThreadExtension() server.Extension {
	return &threadExtension{}
}

func (s *threadExtension) Capabilities(c server.Conn) []string {
	if c.Context().State&imap.AuthenticatedState == 0 {
		return nil
	}

	be, ok := c.Server().Backend.(ThreadBackend)
	if !ok {
		// No backend support, no-op.
		return nil
	}

	var caps []string
	for _, algo := range be.SupportedThreadAlgorithms() {
		caps = append(caps, string("THREAD="+algo))
	}
	return caps
}

func (s *threadExtension) Command(name string) server.HandlerFactory {
	if name == "THREAD" {
		return func() server.Handler {
			return &ThreadHandler{}
		}
	}
	return nil
}
