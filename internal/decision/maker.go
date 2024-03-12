package decision

import (
	"strings"
)

const (
	RangeSeparator = "$%@"
)

type Maker struct {
	stringCaster *StringCaster
	comparer     *Comparer
}

func NewMaker() *Maker {
	return &Maker{
		stringCaster: &StringCaster{},
		comparer: &Comparer{
			int:    &RealFullCompare[int64]{},
			float:  &RealFullCompare[float64]{},
			string: &RealEqualCompare[string]{},
			bool:   &RealEqualCompare[bool]{},
		},
	}
}

func NewMakerForTestConfiguration() *Maker {
	return &Maker{
		stringCaster: &StringCaster{},
		comparer: &Comparer{
			int:    &TestConfigurationFullCompare[int64]{},
			float:  &TestConfigurationFullCompare[float64]{},
			string: &TestConfigurationEqualCompare[string]{},
			bool:   &TestConfigurationEqualCompare[bool]{},
		},
	}
}

func (M *Maker) Decide(schmalue string, value string, compare Compare, parameterType string) bool {
	if compare == Range {
		vals := strings.Split(value, RangeSeparator)
		schmals := strings.Split(schmalue, RangeSeparator)
		result1 := M.Decide(schmals[0], vals[0], LowerEqual, parameterType)
		result2 := M.Decide(schmals[0], vals[1], GreaterEqual, parameterType)
		result3 := M.Decide(schmals[1], vals[0], LowerEqual, parameterType)
		result4 := M.Decide(schmals[1], vals[1], GreaterEqual, parameterType)

		return (result1 == true && result2 == true) || (result3 == true && result4 == true)
	} else {
		switch parameterType {
		case Integer:
			compare1 := M.stringCaster.toInt(schmalue)
			compare2 := M.stringCaster.toInt(value)
			return M.comparer.CompareInt(compare1, compare2, compare)
		case Float:
			compare1 := M.stringCaster.toFloat(schmalue)
			compare2 := M.stringCaster.toFloat(value)
			return M.comparer.CompareFloat(compare1, compare2, compare)
		case DateTime, Date, Time:
			compare1 := M.stringCaster.toInt(schmalue)
			compare2 := M.stringCaster.toInt(value)
			return M.comparer.CompareInt(compare1, compare2, compare)
		case String:
			return M.comparer.CompareString(schmalue, value, compare)
		case Boolean:
			compare1 := M.stringCaster.toBool(schmalue)
			compare2 := M.stringCaster.toBool(value)
			return M.comparer.CompareBool(compare1, compare2, compare)
		default:
			return false
		}
	}
}
