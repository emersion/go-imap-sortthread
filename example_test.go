package sortthread_test

import (
	"log"
	"time"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap-sortthread"
	"github.com/emersion/go-imap/client"
)

func ExampleSortClient() {
	// Assuming c is an IMAP client
	var c *client.Client

	// Create a new SORT client
	sc := sortthread.NewSortClient(c)

	// Check the server supports the extension
	ok, err := sc.SupportSort()
	if err != nil {
		log.Fatal(err)
	} else if !ok {
		log.Fatal("Server doesn't support SORT")
	}

	// Send a SORT command: search for the first 10 messages, sort them by
	// ascending sender and then by descending size
	sortCriteria := []sortthread.SortCriterion{
		{Field: sortthread.SortFrom},
		{Field: sortthread.SortSize, Reverse: true},
	}
	searchCriteria := imap.NewSearchCriteria()
	searchCriteria.SeqNum = new(imap.SeqSet)
	searchCriteria.SeqNum.AddRange(1, 10)
	uids, err := sc.UidSort(sortCriteria, searchCriteria)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(uids)
}

func ExampleThreadClient() {
	// Assuming c is an IMAP client
	var c *client.Client

	// Create a new THREAD client
	sc := sortthread.NewThreadClient(c)

	// Check the server supports the extension
	ok, err := sc.SupportThread()
	if err != nil {
		log.Fatal(err)
	} else if !ok {
		log.Fatal("Server doesn't support THREAD")
	}

	layoutISO := "2006-01-02"
	searchCriteria := imap.NewSearchCriteria()
	date, _ := time.Parse(layoutISO, "2019-07-05")
	searchCriteria.Since = date
	threads, err := sc.UidThread(sortthread.References, searchCriteria)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(threads)
}
