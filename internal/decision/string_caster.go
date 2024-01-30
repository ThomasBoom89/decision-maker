package decision

import "strconv"

type StringCaster struct {
}

func (S *StringCaster) toInt(value string) int64 {
	parsedValue, _ := strconv.ParseInt(value, 10, 64)

	return parsedValue
}

func (S *StringCaster) fromInt(value int64) string {
	return strconv.FormatInt(value, 10)
}
func (S *StringCaster) toFloat(value string) float64 {
	parsedValue, _ := strconv.ParseFloat(value, 64)

	return parsedValue
}

func (S *StringCaster) fromFloat(value float64) string {
	return strconv.FormatFloat(value, 'g', -1, 64)
}
func (S *StringCaster) toDateTime(value string) {
	// todo implement
	//return parsedValue
}
func (S *StringCaster) toBool(value string) bool {
	parsedValue, _ := strconv.ParseBool(value)

	return parsedValue
}
