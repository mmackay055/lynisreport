package lynis

/*
* Author: Matt MacKay
*  Email: mmacKay055@gmail.com
*   Date: 2022-04-06
 */

import (
	"bytes"
	"encoding/json"
	"io"
)

// Process the report by reading from Reader and formatting with the
// OutputFormatter, returns the pointer to the Report struct and a byte array
// of the serialized report
func Process(input io.Reader, output OutputFormatter) (*Report, []byte, error) {
	// CreateReport object from input
	report, err := CreateReport(input)
	if err != nil {
		return nil, nil, err
	}

	return output.Format(report, nil, nil)
}

// OutputFormatter interface for outputing the Report struct to specific format
type OutputFormatter interface {
	// format report into byte slice
	Format(*Report, []byte, error) (*Report, []byte, error)

	// Format to execute after
	Next() OutputFormatter

	// Set formater to execute after
	SetNext(OutputFormatter)
}

// Appends a processor to the end of the processor chain
func SetNext(head OutputFormatter, next OutputFormatter) {
	var of OutputFormatter
	prev := head

	// find last processor
	for of = prev.Next(); of != nil; of = prev.Next() {
		prev = of
	}

	// set processor
	prev.SetNext(next)
}

// OutputFormatter that will format report as a JSON string
type FormatJSON struct {
	next OutputFormatter
}

// Serializes Report into Json byte slice and returns the Report pointer, and
// byte slice
func (fj *FormatJSON) Format(report *Report,
	data []byte, err error) (*Report, []byte, error) {

	// if error exists return the error
	if err != nil {
		return nil, nil, err
	}

	// marshal the Report struct into byte slice
	newdata, err := json.Marshal(report)
	if err != nil {
		return nil, nil, err
	}

	if data == nil {
		data = newdata
	} else {
		// if data already exists in slice
		// append json string and add space between data
		buf := bytes.NewBuffer(data)
		buf.WriteRune(' ')
		buf.Write(newdata)
		data = buf.Bytes()
	}

	// execute the next formatter if it exists
	if fj.Next() != nil {
		return fj.Next().Format(report, data, nil)
	} else {
		return report, data, nil
	}
}

// Returns the next formatter
func (fj *FormatJSON) Next() OutputFormatter {
	return fj.next
}

// Sets the next formatter
func (fj *FormatJSON) SetNext(next OutputFormatter) {
	fj.next = next
}

// OutputFormatter that will format report multiple JSON strings that are
// seperated by newlines. A new JSON string will be generated for each
// Test element that exists in the report. This allows the report to be
// ingested correctly into Elasticsearch that requires flattened objects to
// properly index them
type FormatElasticJSON struct {
	next OutputFormatter
}

// Serializes Report into multiple Json byte slices seperated by new lines and
// returns the Report pointer, and byte slice
func (fj *FormatElasticJSON) Format(report *Report,
	data []byte, err error) (*Report, []byte, error) {

	// return error if it exists
	if err != nil {
		return nil, nil, err
	}

	// serialize TestElements into multiple JSON strings
	newdata, err := report.SerializeForElasticSearch()
	if err != nil {
		return nil, nil, err
	}

	if data == nil {
		data = newdata
	} else {
		// append serialized data if it exists
		buf := bytes.NewBuffer(data)
		buf.WriteRune(' ')
		buf.Write(newdata)
		data = buf.Bytes()
	}

	// execute next formatter if it exists
	if fj.Next() != nil {
		return fj.Next().Format(report, data, nil)
	} else {
		return report, data, nil
	}
}

// Gets next formatter
func (fj *FormatElasticJSON) Next() OutputFormatter {
	return fj.next
}

// Sets next formatter
func (fj *FormatElasticJSON) SetNext(next OutputFormatter) {
	fj.next = next
}

// Formatter that adds timestamp to beginning of serialized data
type FormatTimestamp struct {
	next OutputFormatter
}

// Serializes Report.DateTimeEnd field and adds it to beginning to serialized
// byte slice
func (ft *FormatTimestamp) Format(report *Report,
	data []byte, err error) (*Report, []byte, error) {

        // retun error if it exists
	if err != nil {
		return nil, nil, err
	}

	if data == nil {
		data = []byte(report.DateTimeEnd)
	} else {
                // adds timestamp to beginning of list
		buf := bytes.NewBufferString(report.DateTimeEnd)
		buf.WriteRune(' ')
		buf.Write(data)
		data = buf.Bytes()
	}

        // exexute next formatter
	if ft.Next() != nil {
		return ft.Next().Format(report, data, nil)
	} else {
		return report, data, nil
	}
}

// Get next formatter
func (ft *FormatTimestamp) Next() OutputFormatter {
	return ft.next
}

// Set next formatter
func (ft *FormatTimestamp) SetNext(next OutputFormatter) {
	ft.next = next
}

// Formatter that appends a new line character to end of serialized data slice
type FormatNewLine struct {
	next OutputFormatter
}


// Appends new line character to end of serialized byte slice 
func (fnl *FormatNewLine) Format(report *Report,
	data []byte, err error) (*Report, []byte, error) {
        
        // return error if it exists
	if err != nil {
		return nil, nil, err
	}

	if data == nil {
		data = []byte{'\n'}
	} else {
                // append new line character
		data = append(data, []byte{'\n'}...)
	}

        // exexcute next formatter
	if fnl.Next() != nil {
		return fnl.Next().Format(report, data, nil)
	} else {
		return report, data, nil
	}
}

// Get next formatter
func (fnl *FormatNewLine) Next() OutputFormatter {
	return fnl.next
}

// Set next formatter
func (fnl *FormatNewLine) SetNext(next OutputFormatter) {
	fnl.next = next
}
