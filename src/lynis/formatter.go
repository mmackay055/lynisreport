package lynis

import (
	"encoding/json"
	"io"
        "bytes"
)

func Process(input io.Reader, output OutputFormatter) (*Report, []byte, error) {
	report, err := CreateReport(input)
	if err != nil {
		return nil, nil, err
	}

	return output.Format(report, nil, nil)
}

type OutputFormatter interface {
	Format(*Report, []byte, error) (*Report, []byte, error)
	Next() OutputFormatter
	SetNext(OutputFormatter)
}

func SetNext(head OutputFormatter, next OutputFormatter) {
        var of OutputFormatter
        prev := head
        for of = prev.Next();of != nil; of = prev.Next() {
                prev = of
        }
        prev.SetNext(next)
}

type FormatJSON struct {
	next OutputFormatter
}

func (pj *FormatJSON) Format(report *Report,
	data []byte, err error) (*Report, []byte, error) {
	if err != nil {
		return nil, nil, err
	}

	newdata, err := json.Marshal(report)
	if err != nil {
		return nil, nil, err
	}

	if data == nil {
		data = newdata
	} else {
                // add space between data
                buf := bytes.NewBuffer(data)
                buf.WriteRune(' ')
                buf.Write(newdata)
                data = buf.Bytes()
	}

	if pj.Next() != nil {
		return pj.Next().Format(report, data, nil)
	} else {
		return report, data, nil
	}
}

func (pj *FormatJSON) Next() OutputFormatter {
	return pj.next
}

func (pj *FormatJSON) SetNext(next OutputFormatter) {
	pj.next = next
}

type FormatTimestamp struct {
	next OutputFormatter
}

func (pj *FormatTimestamp) Format(report *Report, data []byte, err error) (*Report, []byte, error) {
	if err != nil {
		return nil, nil, err
	}

	if data == nil {
		data = []byte(report.DateTimeEnd)
	} else {
                buf := bytes.NewBufferString(report.DateTimeEnd)
                buf.WriteRune(' ')
                buf.Write(data)
                data = buf.Bytes()
	}

	if pj.Next() != nil {
		return pj.Next().Format(report, data, nil)
	} else {
		return report, data, nil
	}
}

func (pj *FormatTimestamp) Next() OutputFormatter {
	return pj.next
}

func (pj *FormatTimestamp) SetNext(next OutputFormatter) {
	pj.next = next
}

type FormatNewLine struct {
	next OutputFormatter
}

func (pj *FormatNewLine) Format(report *Report,
	data []byte, err error) (*Report, []byte, error) {
	if err != nil {
		return nil, nil, err
	}

	if data == nil {
		data = []byte{'\n'}
	} else {
                data = append(data, []byte{'\n'}...)
	}

	if pj.Next() != nil {
		return pj.Next().Format(report, data, nil)
	} else {
		return report, data, nil
	}
}

func (pj *FormatNewLine) Next() OutputFormatter {
	return pj.next
}

func (pj *FormatNewLine) SetNext(next OutputFormatter) {
	pj.next = next
}
