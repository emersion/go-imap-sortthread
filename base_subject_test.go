package sortthread

import "testing"

var baseSubjectTests = []struct {
	name       string
	subject    string
	expected   string
	isReplyFwd bool
}{
	{
		name:       "simple",
		subject:    "No Replacement",
		expected:   "No Replacement",
		isReplyFwd: false,
	},
	{
		name:       "reply",
		subject:    "Re: [ocf/puppet] Fix kerberos not booting up correctly [needs testing] (#781)",
		expected:   "[ocf/puppet] Fix kerberos not booting up correctly [needs testing] (#781)",
		isReplyFwd: true,
	},
	{
		name:       "forward",
		subject:    "Fwd: waifus",
		expected:   "waifus",
		isReplyFwd: true,
	},
	{
		name:       "forward_reply",
		subject:    "Fwd: Re: ugh",
		expected:   "ugh",
		isReplyFwd: true,
	},
	{
		name:       "forward_header_simple",
		subject:    "[FWD: simple [extraction]]",
		expected:   "simple [extraction]",
		isReplyFwd: true,
	},
	{
		name:       "foward_header_nested",
		subject:    "Re: [fwd: Re: [OCF] Service update during PG&E outage]",
		expected:   "[OCF] Service update during PG&E outage",
		isReplyFwd: true,
	},
}

func TestBaseSubject(t *testing.T) {
	for _, test := range baseSubjectTests {
		t.Run(test.name, func(t *testing.T) {
			baseSubject, isReplyFwd := GetBaseSubject(test.subject)
			if baseSubject != test.expected {
				t.Errorf("Got %s, Expected %s.", baseSubject, test.expected)
			}
			if !isReplyFwd && test.isReplyFwd {
				t.Errorf("Subject %s should be flagged as reply or forward", test.subject)
			} else if isReplyFwd && !test.isReplyFwd {
				t.Errorf("Subject %s was incorrectly flagged as reply or forward", test.subject)
			}
		})
	}
}
