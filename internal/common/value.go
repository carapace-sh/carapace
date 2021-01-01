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

type ByValue []RawValue

func (a ByValue) Len() int           { return len(a) }
func (a ByValue) Less(i, j int) bool { return a[i].Value < a[j].Value }
func (a ByValue) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
