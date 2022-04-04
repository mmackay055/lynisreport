package lynis

import (
	"errors"
	"strings"
)

//TODO add error type for setFields function to only issue as a warning
type TestElement struct {
	Message  string `json:"message"`
	Details  string `json:"details"`
	Solution string `json:"solution"`
}

func NewTestElement(values []string) (*TestElement, error) {
	length := len(values)
	if length != 3 && length >= 1 {
		// just set message field joined with other fields
		//TODO send warning message as specific error
		return &TestElement{
			Message: strings.Join(values, "|"),
		}, nil
	} else if length != 3 {
		return nil, errors.New("element does not have correct amount of fields")
	}

	return &TestElement{values[0], values[1], values[2]}, nil
}
