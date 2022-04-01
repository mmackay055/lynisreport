package lynis

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
)

const (
	REPORT_NON_LINE_REG string = `(?:^\s*#)|(?:^\s*$)`
)

const (
	KEY_LYNISVER              string = `lynis_version`
	KEY_WARNING               string = `warning[]`
	KEY_SUGGESTION            string = `suggestion[]`
	KEY_REPORT_DATETIME_START string = `report_datetime_start`
	KEY_REPORT_DATETIME_END   string = `report_datetime_end`
)

const (
        VER string = "3.0.7"
)

type Report struct {
	LynisVersion  string           `json:"lynisVersion"`
	DateTimeStart string           `json:"DateTimeStart"`
	DateTimeEnd   string           `json:"DateTimeEnd"`
	Tests         map[string]*Test `json:"tests"`
	nonline       *regexp.Regexp
}

func NewReport() *Report {
	return &Report{
		Tests:   make(map[string]*Test),
		nonline: regexp.MustCompile(REPORT_NON_LINE_REG),
	}

}

func CreateReport(input io.Reader) (*Report, error) {
	report := NewReport()
	err := report.Process(input)
	if err != nil {
		return nil, err
	}
	return report, nil
}

func CheckVersion(ver string) error {
        vCheck := strings.Split(ver, ".")
        vValid := strings.Split(VER, ".")
        if len(vCheck) != 3 {
                return errors.New("Lynis version string invalid format")
        }

        // check all version numbers
        for i, v := range vCheck {
                check,err := strconv.Atoi(v)
                if err != nil {
                        return errors.New("Lynis version string contains non number")
                }
                valid,err := strconv.Atoi(vValid[i])
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

func (r *Report) ProcessLine(line string) error {
	// skip lines commented with '#' or empty
	if r.nonline.Match([]byte(line)) {
		return nil
	}

	key, value, err := parseKeyValue(line)
	if err != nil {
		return nil
	}

	if err := r.Add(key, value); err != nil {
		return err
	}
	return nil
}

func (r *Report) Add(key, value string) error {
	var err error
	switch key {
	case KEY_LYNISVER:
                if err := CheckVersion(value); err != nil {
                        return err
                }
                r.LynisVersion = value
	case KEY_WARNING:
		_, err = r.parseTestValues(value, &Warning{}, AddWarning)
	case KEY_SUGGESTION:
		_, err = r.parseTestValues(value, &Suggestion{}, AddSuggestion)
	case KEY_REPORT_DATETIME_START:
		r.DateTimeStart = value
	case KEY_REPORT_DATETIME_END:
		r.DateTimeEnd = value
	default:
		return nil
	}
	return err
}

func (r *Report) AddTest(name string) *Test {
	test, ok := r.Tests[name]
	if ok {
		return test
	}

	test = NewTest(name)
	r.Tests[name] = test
	return test
}

func (r *Report) parseTestValues(value string, te TestElement, add func(*Test, TestElement)) (*Test, error) {
	values := strings.Split(value, "|")
	values = values[:len(values)-1] // remove last element which is blank

	if len(values) < 2 {
		return nil, errors.New("malformed line no test name or test is missing info")
	}

	// set the fields of the TestElement
	if err := SetFields(values[1:], te); err != nil {
		return nil, err
	}

	// create new test object or get existing
	test := r.AddTest(values[0])

	// add TestElement to Test
	add(test, te)

	return test, nil
}

func parseKeyValue(line string) (key string, value string, err error) {
	keyValue := strings.Split(line, "=")

	if len(keyValue) != 2 {
		return "", "", errors.New("malformed line")
	}

	key = strings.TrimSpace(keyValue[0])
	value = strings.TrimSpace(keyValue[1])

	return
}

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