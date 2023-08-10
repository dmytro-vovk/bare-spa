package types

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNTP_Validate(t *testing.T) {
	testCases := []struct {
		name     string
		receiver NTP
		expected NTP
	}{
		{
			name:     "without servers",
			receiver: NTP{Servers: []string{}},
			expected: NTP{Servers: []string{}},
		},
		{
			name:     "without first server",
			receiver: NTP{Servers: []string{"", "1.lede.pool.ntp.org", "2.lede.pool.ntp.org"}},
			expected: NTP{Servers: []string{"1.lede.pool.ntp.org", "2.lede.pool.ntp.org"}},
		},
		{
			name:     "without middle server",
			receiver: NTP{Servers: []string{"0.lede.pool.ntp.org", "", "2.lede.pool.ntp.org"}},
			expected: NTP{Servers: []string{"0.lede.pool.ntp.org", "2.lede.pool.ntp.org"}},
		},
		{
			name:     "without last server",
			receiver: NTP{Servers: []string{"0.lede.pool.ntp.org", "1.lede.pool.ntp.org", ""}},
			expected: NTP{Servers: []string{"0.lede.pool.ntp.org", "1.lede.pool.ntp.org"}},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.receiver.Validate()
			assert.Equal(t, tc.expected, tc.receiver)
		})
	}
}
