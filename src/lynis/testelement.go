package lynis

import (
	"errors"
)

type TestElement interface {
	setFields(values []string) error
}

type Warning struct {
	Message string `json:"message"`
	Object  string `json:"object"`
	Misc    string `json:"misc"`
}

func (w *Warning) setFields(values []string) error {
	if len(values) != 3 {
		return errors.New("warning is missing fields")
	}
        
	w.Message = values[0]
	w.Object  = values[1]
	w.Misc    = values[2]
        
	return nil
}

type Suggestion struct {
	Message string `json:"message"`
	Misc1   string `json:"misc1"`
	Misc2   string `json:"misc2"`
}

func (s *Suggestion) setFields(values []string) error {
	if len(values) != 3 {
		return errors.New("suggestion is missing fields")
	}

	s.Message = values[0]
	s.Misc1   = values[1]
	s.Misc2   = values[2]

	return nil

}
