package sortthread

import "testing"

var threadTests = []struct {
	name     string
	expected string
	thread   []interface{}
}{
	{
		name:     "simple",
		expected: "(1 2 3 4)",
		thread:   []interface{}{"1", "2", "3", "4"},
	},
	{
		name:     "noparent",
		expected: "(3)(5)",
		thread:   []interface{}{[]interface{}{"3"}, []interface{}{"5"}},
	},
	{
		name:     "nested",
		expected: "(4 5 (6) (7 8))",
		thread:   []interface{}{"4", "5", []interface{}{"6"}, []interface{}{"7", "8"}},
	},
}

func TestThreadParsing(t *testing.T) {
	for _, test := range threadTests {
		t.Run(test.name, func(t *testing.T) {
			threads, err := parseThreadResp(test.thread)
			if err != nil {
				t.Error("Expected no error while parsing thread but got:", err)
			}
			var s string
			for _, t := range threads {
				s += t.String()
			}
			if s != test.expected {
				t.Errorf("Got %s; Expected %s", s, test.expected)
			}
		})
	}
}
