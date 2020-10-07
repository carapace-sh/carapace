package common

type Candidate struct {
	Value       string
	Display     string
	Description string
}

type Candidates []Candidate

func CandidateFromValues(values ...string) Candidates {
	candidates := make([]Candidate, len(values))
	for index, val := range values {
		candidates[index] = Candidate{Value: val, Display: val}
	}
	return candidates
}
