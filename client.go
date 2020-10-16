package sortthread

import (
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-imap/commands"
)

// SortClient is a SORT client.
type SortClient struct {
	c *client.Client
}

// ThreadClient is a THREAD client.
type ThreadClient struct {
	c *client.Client
}

// NewClient creates a new SORT client.
func NewSortClient(c *client.Client) *SortClient {
	return &SortClient{c: c}
}

// SupportSort returns true if the remote server supports the extension.
func (c *SortClient) SupportSort() (bool, error) {
	return c.c.Support(SortCapability)
}

func (c *SortClient) sort(uid bool, sortCriteria []SortCriterion, searchCriteria *imap.SearchCriteria) ([]uint32, error) {
	if c.c.State() != imap.SelectedState {
		return nil, client.ErrNoMailboxSelected
	}

	var cmd imap.Commander
	cmd = &SortCommand{
		SortCriteria:   sortCriteria,
		Charset:        "UTF-8",
		SearchCriteria: searchCriteria,
	}
	if uid {
		cmd = &commands.Uid{Cmd: cmd}
	}

	res := new(SortResponse)

	status, err := c.c.Execute(cmd, res)
	if err != nil {
		return nil, err
	}

	return res.Ids, status.Err()
}

func (c *SortClient) Sort(sortCriteria []SortCriterion, searchCriteria *imap.SearchCriteria) ([]uint32, error) {
	return c.sort(false, sortCriteria, searchCriteria)
}

func (c *SortClient) UidSort(sortCriteria []SortCriterion, searchCriteria *imap.SearchCriteria) ([]uint32, error) {
	return c.sort(true, sortCriteria, searchCriteria)
}

// NewClient creates a new THREAD client
func NewThreadClient(c *client.Client) *ThreadClient {
	return &ThreadClient{c: c}
}

// SupportThread returns true if the remote server supports the extension.
func (c *ThreadClient) SupportThread() (bool, error) {
	for _, capability := range ThreadCapabilities {
		ok, err := c.c.Support(capability)
		if err != nil {
			return false, err
		} else if ok {
			return true, nil
		}
	}
	return false, nil
}

func (c *ThreadClient) thread(uid bool, algorithm ThreadAlgorithm, searchCriteria *imap.SearchCriteria) ([]*Thread, error) {
	if c.c.State() != imap.SelectedState {
		return nil, client.ErrNoMailboxSelected
	}

	var cmd imap.Commander
	cmd = &ThreadCommand{
		Algorithm:      algorithm,
		Charset:        "UTF-8",
		SearchCriteria: searchCriteria,
	}

	if uid {
		cmd = &commands.Uid{Cmd: cmd}
	}

	res := new(ThreadResponse)

	status, err := c.c.Execute(cmd, res)
	if err != nil {
		return nil, err
	}

	return res.Threads, status.Err()
}

func (c *ThreadClient) Thread(algorithm ThreadAlgorithm, searchCriteria *imap.SearchCriteria) ([]*Thread, error) {
	return c.thread(false, algorithm, searchCriteria)
}

func (c *ThreadClient) UidThread(algorithm ThreadAlgorithm, searchCriteria *imap.SearchCriteria) ([]*Thread, error) {
	return c.thread(true, algorithm, searchCriteria)
}
