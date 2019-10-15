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
}

func TestBaseSubject(t *testing.T) {
	for _, test := range baseSubjectTests {
		t.Run(test.name, func(t *testing.T) {
			var isReplyFwd bool
			baseSubject, err := GetBaseSubject(test.subject, &isReplyFwd)
			if err != nil {
				t.Error("Expected no error while parsing subject but got:", err)
			}
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
