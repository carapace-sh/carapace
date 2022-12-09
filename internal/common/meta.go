package common

type Meta struct {
	Messages Messages
	Nospace  SuffixMatcher
	Usage    string
}

func (m *Meta) Merge(other Meta) {
	if other.Usage != "" {
		m.Usage = other.Usage
	}
	m.Nospace.Merge(other.Nospace)
	m.Messages.Merge(other.Messages)
}
