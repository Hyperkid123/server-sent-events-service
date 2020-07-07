package conditions

import (
	"fmt"
)

// ApprovalCondition is used to determine if the event can be emmited
func ApprovalCondition(msg string, accountNumber string) bool {

	fmt.Println("Got new approval event", msg)
	fmt.Println("Using account", accountNumber)

	return true
}
