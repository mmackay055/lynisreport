package lynis

type Test struct {
	Name        string        `json:"test"`
	Warnings    []*Warning    `json:"warnings"`
	Suggestions []*Suggestion `json:"suggestions"`
}

func NewTest(name string) *Test {
	return &Test{
		Name:        name,
		Warnings:    make([]*Warning, 0),
		Suggestions: make([]*Suggestion, 0),
	}
}

func AddWarning(t *Test, te TestElement) {
	t.Warnings = append(t.Warnings, te.(*Warning))
}

func AddSuggestion(t *Test, te TestElement) {
	t.Suggestions = append(t.Suggestions, te.(*Suggestion))
}
