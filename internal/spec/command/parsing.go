package command

type Parsing string

const (
	DEFAULT          Parsing = ""                 // INTERSPERSED but allows implicit changes
	INTERSPERSED     Parsing = "interspersed"     // mixed flags and positional arguments
	NON_INTERSPERSED Parsing = "non-interspersed" // flag parsing stopped after first positional argument
	DISABLED         Parsing = "disabled"         // flag parsing disabled
)
