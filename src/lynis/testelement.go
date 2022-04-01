package lynis

import (
	"errors"
	"strings"
        "reflect"
)

//TODO add error type for setFields function to only issue as a warning

type TestElement interface {
	setFields(values []string) error
}

func SetFields(values []string, te TestElement) error {
        length := len(values)
	if length != 3 && length >= 1 {
                // just set first field joined with other fields
                reflect.ValueOf(te).Elem().Field(0).SetString(strings.Join(values, "|"))

                //TODO send warning message as specific error
		return nil
	} else if length != 3 {
                return errors.New("element does not have correct amount of fields")
        }

        // set TestElement fields
        return te.setFields(values) 
}
type Warning struct {
	Message string `json:"message"`
	Object  string `json:"object"`
	Misc    string `json:"misc"`
}


func (w *Warning) setFields(values []string) error {
        w.Message = values[0]
        w.Object = values[1]
        w.Misc = values[2]
	return nil
}

type Suggestion struct {
	Message  string `json:"message"`
	Details  string `json:"details"`
	Solution string `json:"solution"`
}

func (s *Suggestion) setFields(values []string) error {
        s.Message = values[0]
        s.Details = values[1]
        s.Solution = values[2]
	return nil
}
