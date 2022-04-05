package lynis

type Test struct {
	Name        string         `json:"testname"`
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

func (t *Test) CreateTestElementElastics(r *Report) []*TestElementElastic {
        tees := make([]*TestElementElastic,
        len(t.Warnings) + len(t.Suggestions))

        j := 0
        for i, w := range t.Warnings {
                tees[i],_ = CreateTestElementElastic("warning", r, t, w)
                j = i + 1
        }

        for i, s := range t.Suggestions {
                tees[i + j], _ = CreateTestElementElastic("suggestion", r, t, s)
        }

        return tees
}

