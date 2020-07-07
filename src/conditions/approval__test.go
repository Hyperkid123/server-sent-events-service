package conditions

import "testing"

func TestApprovalCondition(t *testing.T) {
	got := ApprovalCondition("Event message", "123")
	if got != true {
		t.Errorf("ApprovalCondition('Event message', '123') = %t; want true", got)
	}
}
