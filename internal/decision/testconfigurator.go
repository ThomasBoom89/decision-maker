package decision

import (
	"math"
)

type TestConfigurator struct {
	stringCaster *StringCaster
}

type ValueTypeComparer struct {
	Name     string
	Value    string
	Type     string
	Comparer Compare
}

func NewTestConfigurator() *TestConfigurator {
	return &TestConfigurator{}
}

func (T *TestConfigurator) Create(valueTypeComparers map[uint]ValueTypeComparer) map[string]string {
	ret := make(map[string]string)
	for _, valueTypeComparer := range valueTypeComparers {
		var value string
		switch valueTypeComparer.Type {
		case Boolean:
			value = T.getForBool(valueTypeComparer.Value, valueTypeComparer.Comparer)
		case Integer:
			value = T.getForInteger(valueTypeComparer.Value, valueTypeComparer.Comparer)
		case String:
			value = T.getForString(valueTypeComparer.Value, valueTypeComparer.Comparer)
		case Float:
			value = T.getForFloat(valueTypeComparer.Value, valueTypeComparer.Comparer)
		case DateTime:
		// todo
		default:
			panic("nope")
		}
		ret[valueTypeComparer.Name] = value
	}

	return ret
}

func (T *TestConfigurator) getForBool(value string, comparer Compare) string {
	switch comparer {
	case Equal, NotEqual:
		return value
	default:
		panic("shit21")
	}
}

func (T *TestConfigurator) getForInteger(value string, comparer Compare) string {
	switch comparer {
	case Equal, LowerEqual, GreaterEqual, NotEqual:
		return value
	case GreaterThan:
		result := T.stringCaster.toInt(value)
		result++
		return T.stringCaster.fromInt(result)
	case LowerThan:
		result := T.stringCaster.toInt(value)
		result--
		return T.stringCaster.fromInt(result)
	default:
		panic("darf nicht passieren")
	}
}

func (T *TestConfigurator) getForString(value string, comparer Compare) string {
	switch comparer {
	case Equal, NotEqual:
		return value
	default:
		panic("darf nicht passieren")
	}
}

func (T *TestConfigurator) getForFloat(value string, comparer Compare) string {
	switch comparer {
	case Equal, LowerEqual, GreaterEqual, NotEqual:
		return value
	case GreaterThan:
		result := T.stringCaster.toFloat(value)
		next := math.Nextafter(result, math.MaxFloat64)
		return T.stringCaster.fromFloat(next)
	case LowerThan:
		result := T.stringCaster.toFloat(value)
		next := math.Nextafter(result, math.MaxFloat64)
		return T.stringCaster.fromFloat(result - (next - result))
	default:
		panic("darf nicht passieren")
	}
}
