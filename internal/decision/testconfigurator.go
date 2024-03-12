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
		case DateTime, Date, Time:
			value = T.getForInteger(valueTypeComparer.Value, valueTypeComparer.Comparer)
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
	case Equal, NotEqual, Range:
		return value
	case GreaterThan, GreaterEqual:
		result := int64(math.MaxInt64)
		return T.stringCaster.fromInt(result)
	case LowerThan, LowerEqual:
		result := int64(math.MinInt64)
		return T.stringCaster.fromInt(result)
	default:
		panic("darf nicht passieren integer")
	}
}

func (T *TestConfigurator) getForString(value string, comparer Compare) string {
	switch comparer {
	case Equal, NotEqual:
		return value
	default:
		panic("darf nicht passieren string")
	}
}

func (T *TestConfigurator) getForFloat(value string, comparer Compare) string {
	switch comparer {
	case Equal, NotEqual, Range:
		return value
	case GreaterThan, GreaterEqual:
		next := math.MaxFloat64
		return T.stringCaster.fromFloat(next)
	case LowerThan, LowerEqual:
		next := math.MaxFloat64 * -1
		return T.stringCaster.fromFloat(next)
	default:
		panic("darf nicht passieren float")
	}
}
