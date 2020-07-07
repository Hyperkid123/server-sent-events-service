package mutators

import (
	"encoding/json"
	"fmt"
)

type ApprovalMessage struct {
	RequestID string `json:"request_id,omitempty"`
	GroupName string `json:"group_name,omitempty"`
	Decision  string `json:"decision,omitempty"`
	Reason    string `json:"reason,omitempty"`
	Mutation  bool   `json:"mutation,omitempty"`
}

func ApprovalMutator(msg []byte) []byte {
	var messageObject ApprovalMessage
	err := json.Unmarshal(msg, &messageObject)
	if err != nil {
		fmt.Println("Unable to decode approval message", err)
		return msg
	}
	messageObject.Mutation = true

	response, err := json.Marshal(messageObject)
	if err != nil {
		fmt.Println("Unable to mutate approval message", err)
		return msg
	}
	return response
}
