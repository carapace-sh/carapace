package command

type Flag struct {
	Longhand   string
	Shorthand  string
	Usage      string
	Repeatable bool
	Optarg     bool
	Value      bool
	Hidden     bool
	Required   bool
	Persistent bool
}

func (f Flag) format() string {
	var s string

	if f.Shorthand != "" {
		s += f.Shorthand
		if f.Longhand != "" {
			s += ", "
		}
	}

	if f.Longhand != "" {
		s += f.Longhand
	}

	switch {
	case f.Optarg:
		s += "?"
	case f.Value:
		s += "="
	}

	if f.Repeatable {
		s += "*"
	}

	if f.Required {
		s += "!"
	}

	if f.Hidden {
		s += "&"
	}

	return s
}
