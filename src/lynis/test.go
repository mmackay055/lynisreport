package lynis

type Test struct {
	Name        string         `json:"name"`
	Warnings    []*TestElement `json:"warnings"`
	Suggestions []*TestElement `json:"suggestions"`
}

func NewTest(name string) *Test {
	return &Test{
		Name:        name,
		Warnings:    make([]*TestElement, 0),
		Suggestions: make([]*TestElement, 0),
	}
}

func AddWarning(t *Test, te *TestElement) {
	t.Warnings = append(t.Warnings, te)
}

func AddSuggestion(t *Test, te *TestElement) {
	t.Suggestions = append(t.Suggestions, te)
}
