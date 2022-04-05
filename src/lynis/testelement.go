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

type TestElementElastic struct {
        Name string `json:"name"`
	Type          string `json:"type"`
	LynisVersion  string `json:"lynisVersion"`
	DateTimeStart string `json:"datetime_start"`
	DateTimeEnd   string `json:"datetime_end"`
	Message       string `json:"message"`
	Details       string `json:"details"`
	Solution      string `json:"solution"`
}

// used to create a flattened object to ingest into Elasiticsearch
func CreateTestElementElastic(typ string, r *Report, t *Test, te *TestElement) (*TestElementElastic, error) {
        return &TestElementElastic{ 
                Name:t.Name,
                Type: typ,
                LynisVersion: r.LynisVersion,
                DateTimeStart: r.DateTimeStart,
                DateTimeEnd: r.DateTimeEnd,
                Message: te.Message,
                Details: te.Details,
                Solution: te.Solution,
        }, nil
}
