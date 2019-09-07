package sortthread

import (
	"reflect"
	"testing"
)

var threadTests = []struct {
	name     string
	str      string
	expected []*Thread
	thread   []interface{}
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
		thread: []interface{}{"1", "2", "3", "4"},
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
		thread: []interface{}{[]interface{}{"3"}, []interface{}{"5"}},
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
		thread: []interface{}{"4", "5", []interface{}{"6"}, []interface{}{"7", "8"}},
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
		thread: []interface{}{
			[]interface{}{"2"},
			[]interface{}{"3", "6",
				[]interface{}{
					"4", "23",
				},
				[]interface{}{
					"44", "7", "96",
				},
			},
		},
	},
}

func TestThreadParsing(t *testing.T) {
	for _, test := range threadTests {
		t.Run(test.name, func(t *testing.T) {
			threads, err := parseThreadResp(test.thread)
			if err != nil {
				t.Error("Expected no error while parsing thread but got:", err)
			}
			if !reflect.DeepEqual(test.expected, threads) {
				t.Errorf("Could not parse %s", test.str)
			}
		})
	}
}
