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

func (fj *FormatJSON) Format(report *Report,
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

	if fj.Next() != nil {
		return fj.Next().Format(report, data, nil)
	} else {
		return report, data, nil
	}
}

func (fj *FormatJSON) Next() OutputFormatter {
	return fj.next
}

func (fj *FormatJSON) SetNext(next OutputFormatter) {
	fj.next = next
}

type FormatTimestamp struct {
	next OutputFormatter
}

func (ft *FormatTimestamp) Format(report *Report, data []byte, err error) (*Report, []byte, error) {
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

	if ft.Next() != nil {
		return ft.Next().Format(report, data, nil)
	} else {
		return report, data, nil
	}
}

func (ft *FormatTimestamp) Next() OutputFormatter {
	return ft.next
}

func (ft *FormatTimestamp) SetNext(next OutputFormatter) {
	ft.next = next
}

type FormatNewLine struct {
	next OutputFormatter
}

func (fnl *FormatNewLine) Format(report *Report,
	data []byte, err error) (*Report, []byte, error) {
	if err != nil {
		return nil, nil, err
	}

	if data == nil {
		data = []byte{'\n'}
	} else {
                data = append(data, []byte{'\n'}...)
	}

	if fnl.Next() != nil {
		return fnl.Next().Format(report, data, nil)
	} else {
		return report, data, nil
	}
}

func (fnl *FormatNewLine) Next() OutputFormatter {
	return fnl.next
}

func (fnl *FormatNewLine) SetNext(next OutputFormatter) {
	fnl.next = next
}
