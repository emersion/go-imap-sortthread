package sortthread

import "testing"

var baseSubjectTests = []struct {
	name     string
	subject  string
	expected string
}{
	{
		name:     "simple",
		subject:  "No Replacement",
		expected: "No Replacement",
	},
	{
		name:     "reply",
		subject:  "Re: [ocf/puppet] Fix kerberos not booting up correctly [needs testing] (#781)",
		expected: "[ocf/puppet] Fix kerberos not booting up correctly [needs testing] (#781)",
	},
	{
		name:     "forward",
		subject:  "Fwd: waifus",
		expected: "waifus",
	},
	{
		name:     "forward_reply",
		subject:  "Fwd: Re: ugh",
		expected: "ugh",
	},
}

func TestBaseSubject(t *testing.T) {
	for _, test := range baseSubjectTests {
		t.Run(test.name, func(t *testing.T) {
			baseSubject, err := GetBaseSubject(test.subject)
			if err != nil {
				t.Error("Expected no error while parsing subject but got:", err)
			}
			if baseSubject != test.expected {
				t.Errorf("Got %s, Expected %s.", baseSubject, test.expected)
			}
		})
	}
}
