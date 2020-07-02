package enhancers

import "testing"

func TestApprovalEnhancer(t *testing.T) {
	got := ApprovalEnhancer("Event message", "123")
	if got != true {
		t.Errorf("ApprovalEnhancer('Event message', '123') = %t; want true", got)
	}	
}