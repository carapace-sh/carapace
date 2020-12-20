package common

type RawValue struct {
	Value       string
	Display     string
	Description string
}

func RawValuesFrom(values ...string) []RawValue {
	rawValues := make([]RawValue, len(values))
	for index, val := range values {
		rawValues[index] = RawValue{Value: val, Display: val}
	}
	return rawValues
}
