package mutators

import (
	"fmt"
	"testing"
)

func TestApprovalMutator(t *testing.T) {
	var tests = []struct {
		in  string
		out string
	}{
		{"{\"request_id\":\"1\",\"group_name\":\"group 1\"}", "{\"request_id\":\"1\",\"group_name\":\"group 1\",\"mutation\":true}"}, // adds mutation
		{"{\"unknown\":\"1\",\"group_name\":\"group 1\"}", "{\"group_name\":\"group 1\",\"mutation\":true}"},                         // omits unknown keys
		{"invalid json", "invalid json"}, // returns original message when json operation fails
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("%s", string(tt.in))
		t.Run(testname, func(t *testing.T) {
			result := string(ApprovalMutator([]byte(tt.in)))
			expected := tt.out
			if result != expected {
				t.Errorf("got %s, want %s", result, expected)
			}
		})
	}
}
