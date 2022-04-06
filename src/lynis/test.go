package lynis

/*
* Author: Matt MacKay
*  Email: mmacKay055@gmail.com
*   Date: 2022-04-06
 */

// Test struct that represents a test performed in Lynis scan
type Test struct {
	Name        string         `json:"testname"`
	Warnings    []*TestElement `json:"warnings"`
	Suggestions []*TestElement `json:"suggestions"`
	report      *Report
}

// Create new Test object
func NewTest(name string, r *Report) *Test {
	return &Test{
		Name:        name,
		Warnings:    make([]*TestElement, 0),
		Suggestions: make([]*TestElement, 0),
		report:      r,
	}
}

// Adds TestElement to Warnings map
func AddWarning(t *Test, te *TestElement) {
	t.Warnings = append(t.Warnings, te)
}

// Adds TestElement to Suggestions map
func AddSuggestion(t *Test, te *TestElement) {
	t.Suggestions = append(t.Suggestions, te)
}

// Creates TestElementElastic elements from test and returns them as a slice
func (t *Test) CreateTestElementElastics() []*TestElementElastic {
        // Create slice to fit all Warnings and suggestions
	tees := make([]*TestElementElastic,
		len(t.Warnings)+len(t.Suggestions))

	j := 0

        // Add warnings to slice 
	for i, w := range t.Warnings {
		tees[i], _ = CreateTestElementElastic("warning",
			t.report, t, w)
		j = i + 1
	}

        // add suggestions to slice
	for i, s := range t.Suggestions {
		tees[i+j], _ = CreateTestElementElastic("suggestion",
			t.report, t, s)
	}

        // return slice
	return tees
}
