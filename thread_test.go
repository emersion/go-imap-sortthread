package sortthread

import (
	"reflect"
	"strings"
	"testing"

	"github.com/emersion/go-imap"
)

var threadTests = []struct {
	name     string
	str      string
	expected []*Thread
	response []interface{}
}{
	{
		name: "simple",
		str:  "(1 2 3 4)",
		expected: []*Thread{
			&Thread{
				Id: 1,
				Children: []*Thread{
					&Thread{
						Id: 2,
						Children: []*Thread{
							&Thread{
								Id: 3,
								Children: []*Thread{
									&Thread{
										Id:       4,
										Children: nil,
									},
								},
							},
						},
					},
				},
			},
		},
		response: []interface{}{[]interface{}{imap.RawString("1"), imap.RawString("2"), imap.RawString("3"), imap.RawString("4")}},
	},
	{
		name: "noparent",
		str:  "(3 5)",
		expected: []*Thread{
			&Thread{
				Id:       3,
				Children: nil,
			},
			&Thread{
				Id:       5,
				Children: nil,
			},
		},
		response: []interface{}{[]interface{}{imap.RawString("3")}, []interface{}{imap.RawString("5")}},
	},
	{
		name: "nested",
		str:  "(4 5 (6) (7 8))",
		expected: []*Thread{
			&Thread{
				Id: 4,
				Children: []*Thread{
					&Thread{
						Id: 5,
						Children: []*Thread{
							&Thread{
								Id:       6,
								Children: nil,
							},
							&Thread{
								Id: 7,
								Children: []*Thread{
									&Thread{
										Id:       8,
										Children: nil,
									},
								},
							},
						},
					},
				},
			},
		},
		response: []interface{}{[]interface{}{
			imap.RawString("4"),
			imap.RawString("5"),
			[]interface{}{imap.RawString("6")},
			[]interface{}{imap.RawString("7"), imap.RawString("8")},
		}},
	},
	{
		name: "rfc",
		str:  "(2)(3 6 (4 23)(44 7 96))",
		expected: []*Thread{
			&Thread{
				Id:       2,
				Children: nil,
			},
			&Thread{
				Id: 3,
				Children: []*Thread{
					&Thread{
						Id: 6,
						Children: []*Thread{
							&Thread{
								Id: 4,
								Children: []*Thread{
									&Thread{
										Id:       23,
										Children: nil,
									},
								},
							},
							&Thread{
								Id: 44,
								Children: []*Thread{
									&Thread{
										Id: 7,
										Children: []*Thread{
											&Thread{
												Id:       96,
												Children: nil,
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		response: []interface{}{
			[]interface{}{imap.RawString("2")},
			[]interface{}{imap.RawString("3"), imap.RawString("6"),
				[]interface{}{
					imap.RawString("4"), imap.RawString("23"),
				},
				[]interface{}{
					imap.RawString("44"), imap.RawString("7"), imap.RawString("96"),
				},
			},
		},
	},
}

func TestThreadParsing(t *testing.T) {
	for _, test := range threadTests {
		t.Run(test.name, func(t *testing.T) {
			threads, err := parseThreadResp(test.response)
			if err != nil {
				t.Error("Expected no error while parsing thread but got:", err)
			}
			if !reflect.DeepEqual(test.expected, threads) {
				t.Errorf("Could not parse %s", test.str)
			}
		})
	}
}

func TestThreadFormatting(t *testing.T) {
	for _, test := range threadTests {
		// noparent case has destructive parsing - parser disregards 'null' parent
		// and resulting tree corresponds to a different response.
		if test.name == "noparent" {
			continue
		}

		t.Run(test.name, func(t *testing.T) {
			fields := formatThreadResp(test.expected)
			if !reflect.DeepEqual(fields, test.response) {
				t.Errorf("Could not format %s properly", test.str)
				t.Logf("Want: %#+v", test.response)
				t.Logf("Got:  %#+v", fields)
			}
		})
	}
}
