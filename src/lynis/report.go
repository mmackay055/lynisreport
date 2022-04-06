package lynis

/*
* Author: Matt MacKay
*  Email: mmacKay055@gmail.com
*   Date: 2022-04-06
 */

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Regex format string to find lines in Lynis report that don't need to be
// processed in report
const (
	REPORT_NON_LINE_REG string = `(?:^\s*#)|(?:^\s*$)`
)

// Key definitions of elements in Lynis report
const (
	// Version of lynis
	KEY_LYNISVER string = `lynis_version`

	// Lynis warning
	KEY_WARNING string = `warning[]`

	// Lynis suggestion
	KEY_SUGGESTION string = `suggestion[]`

	// Lynis start time
	KEY_REPORT_DATETIME_START string = `report_datetime_start`

	// Lynis end time
	KEY_REPORT_DATETIME_END string = `report_datetime_end`
)

const (
	// Minimum compatable version of Lynis
	VER string = "3.0.7"
)

// Report struct that represents a Lynis Report
type Report struct {
	LynisVersion  string           `json:"lynisVersion"`
	DateTimeStart string           `json:"datetime_start"`
	DateTimeEnd   string           `json:"datetime_end"`
	Tests         map[string]*Test `json:"tests"`
	nonline       *regexp.Regexp   // regex used to determine non elements
}

// Initializes a new report
func NewReport() *Report {
	return &Report{
		Tests:   make(map[string]*Test),
		nonline: regexp.MustCompile(REPORT_NON_LINE_REG),
	}
}

// Creates report from Reader, returns Report pointer that is created
func CreateReport(input io.Reader) (*Report, error) {
	// initialize report
	report := NewReport()

	// read report and parse it
	err := report.Process(input)
	if err != nil {
		return nil, err
	}
	return report, nil
}

// Checks that the version of Lynis is compatable
func CheckVersion(ver string) error {
	// parse version string
	vCheck := strings.Split(ver, ".")

	// parse valid version string
	vValid := strings.Split(VER, ".")
	if len(vCheck) != 3 {
		return errors.New("Lynis version string invalid format")
	}

	// check version is greater or equal to valid version
	for i, v := range vCheck {
		check, err := strconv.Atoi(v)
		if err != nil {
			return errors.New("Lynis version string contains non number")
		}

                // panic since internal version string is invalid
		valid, err := strconv.Atoi(vValid[i])
		if err != nil {
			panic(err)
		}
		if check < valid {
			return errors.New("Lynis version is not compatable")
		} else if check > valid {
			break // version is greater than valid version
		}
		// continue checking if version number is equal
	}

	return nil
}

// Process line from Lynis report
func (r *Report) ProcessLine(line string) error {
	// skip lines commented with '#' or empty
	if r.nonline.Match([]byte(line)) {
		return nil
	}

	// get key and value from line
	key, value, err := parseKeyValue(line)
	if err != nil {
                // ignore line
                // TODO send warning message
		return nil
	}

	// add key and value to report
	if err := r.Add(key, value); err != nil {
		return err
	}
	return nil
}

// Adds key and value to the Report struct, if key is not supported it is
// ignored
func (r *Report) Add(key, value string) error {
	var err error

	// process key
	switch key {
	case KEY_LYNISVER:
		// check lynis version
		err = CheckVersion(value)
		r.LynisVersion = value
	case KEY_WARNING:
		// add value to Warning slice for test
		_, err = r.parseTestValues(value, AddWarning)
	case KEY_SUGGESTION:
		// add value to Suggestion slice for test
		_, err = r.parseTestValues(value, AddSuggestion)
	case KEY_REPORT_DATETIME_START:
		// set date time start
		r.DateTimeStart, err = FormatTime(value)
	case KEY_REPORT_DATETIME_END:
		// set date time end
		r.DateTimeEnd, err = FormatTime(value)
	default:
		return nil
	}
	return err
}

// Serialize Report struct so it is compatable to be ingested by Elasticsearch
func (r *Report) SerializeForElasticSearch() ([]byte, error) {
	teesData := make([]byte, 0)
	tees, _ := r.CreateTestElementElastics()

	// Marshal each test specific for Elasticsearch
	for _, te := range tees {

		data, err := json.Marshal(te)
		if err != nil {
			return nil, err
		}

		teesData = append(teesData, data...)
		teesData = append(teesData, '\n')
	}

	return teesData, nil
}

// Creates a slice of TestElementElastic elements
func (r *Report) CreateTestElementElastics() ([]*TestElementElastic, error) {

	tees := make([]*TestElementElastic, 0)

	for _, t := range r.Tests {
		tees = append(tees, t.CreateTestElementElastics()...)
	}

	return tees, nil
}

// Formats the Lynis time fields to ISO8601
func FormatTime(timestr string) (string, error) {
	// get current time
	now := time.Now()

	// add time zone info to date string
	timefmt, err := time.Parse("2006-01-02 15:04:05-0700",
		timestr+now.Format("-0700"))

	// return formatted time
	return timefmt.Format("2006-01-02T15:04:05-0700"), err
}

// Adds test to report by creating a new Test with the name passed to function
// or returns already existing Test
func (r *Report) AddTest(name string) *Test {
	test, ok := r.Tests[name]
	// check if test already exists
	if ok {
		return test
	}

	test = NewTest(name, r)
	r.Tests[name] = test
	return test
}

// Parses the value string retrieved from report and adds to test with add
// function provided. Values are expected to be separated by |
func (r *Report) parseTestValues(value string,
	add func(*Test, *TestElement)) (*Test, error) {

	// split values
	values := strings.Split(value, "|")
	values = values[:len(values)-1] // remove last element which is blank

	// if test name is not provided return error
	if len(values) < 2 {
		return nil, errors.New("malformed line no test name or test is missing info")
	}

	// set the fields of the TestElement
	te, err := NewTestElement(values[1:])
	if err != nil {
		return nil, err
	}

	// create new test object or get existing
	test := r.AddTest(values[0])

	// add TestElement to Test
	add(test, te)

	return test, nil
}

// Parses key value pair from Lynis report
func parseKeyValue(line string) (key string, value string, err error) {
	keyValues := strings.SplitN(line, "=",2)
	if len(keyValues) != 2 {
		return "", "", errors.New("malformed line")
	}

	key = strings.TrimSpace(keyValues[0])
	value = strings.TrimSpace(keyValues[1])

	return
}

// Process the Lynis report from the Reader
func (r *Report) Process(input io.Reader) error {

	bufinput := bufio.NewReader(input)

	// process all lines in file
	for i := 1; ; i++ {
		line, readerr := bufinput.ReadString('\n')
		if readerr != nil && readerr != io.EOF {
			return readerr
		}

		if err := r.ProcessLine(line); err != nil {
			return errors.New(fmt.Sprintf("line %d: %s", i, err))
		}

		if readerr == io.EOF {
			break
		}
	}
	return nil
}
