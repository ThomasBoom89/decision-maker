package decision

type Compare string

const (
	GreaterThan  Compare = "gt"
	GreaterEqual Compare = "ge"
	LowerThan    Compare = "lt"
	LowerEqual   Compare = "le"
	Equal        Compare = "eq"
	NotEqual     Compare = "ne"
	Range        Compare = "range"
)

const (
	Float    = "float"
	Integer  = "int"
	DateTime = "datetime"
	String   = "string"
	Boolean  = "bool"
	Date     = "date"
	Time     = "time"
)

type Comparer struct {
	int    FullCompare[int64]
	float  FullCompare[float64]
	string EqualCompare[string]
	bool   EqualCompare[bool]
}

func (C *Comparer) CompareInt(schmalue int64, value int64, compare Compare) bool {
	return C.int.CompareFull(schmalue, value, compare)
}

func (C *Comparer) CompareFloat(schmalue float64, value float64, compare Compare) bool {
	return C.float.CompareFull(schmalue, value, compare)
}

func (C *Comparer) CompareString(schmalue string, value string, compare Compare) bool {
	return C.string.CompareEqual(schmalue, value, compare)
}

func (C *Comparer) CompareBool(schmalue bool, value bool, compare Compare) bool {
	return C.bool.CompareEqual(schmalue, value, compare)
}

func (C *Comparer) CompareDateTime(schmalue int64, value int64, compare Compare) bool {
	return C.int.CompareFull(schmalue, value, compare)
}

type FullCompare[T int64 | float64] interface {
	CompareFull(schmalue T, value T, compare Compare) bool
}

type EqualCompare[T string | bool] interface {
	CompareEqual(schmalue T, value T, compare Compare) bool
}

type RealFullCompare[T int64 | float64] struct {
}
type RealEqualCompare[T string | bool] struct {
}

func (F *RealFullCompare[T]) CompareFull(schmalue T, value T, compare Compare) bool {
	switch compare {
	case GreaterThan:
		return schmalue < value
	case GreaterEqual:
		return schmalue <= value
	case LowerThan:
		return schmalue > value
	case LowerEqual:
		return schmalue >= value
	case Equal:
		return schmalue == value
	case NotEqual:
		return schmalue != value
	default:
		return false
	}
}

func (E *RealEqualCompare[T]) CompareEqual(schmalue T, value T, compare Compare) bool {
	switch compare {
	case Equal:
		return schmalue == value
	case NotEqual:
		return schmalue != value
	default:
		return false
	}
}

type TestConfigurationFullCompare[T int64 | float64] struct {
}
type TestConfigurationEqualCompare[T string | bool] struct {
}

func (I *TestConfigurationFullCompare[T]) CompareFull(schmalue T, value T, compare Compare) bool {
	switch compare {
	case GreaterThan:
		return schmalue < value
	case GreaterEqual:
		return schmalue <= value
	case LowerThan:
		return schmalue > value
	case LowerEqual:
		return schmalue >= value
	case Equal, NotEqual:
		return schmalue == value
	default:
		return false
	}
}

func (D *TestConfigurationEqualCompare[T]) CompareEqual(schmalue T, value T, compare Compare) bool {
	switch compare {
	case Equal, NotEqual:
		return schmalue == value
	default:
		return false
	}
}

func GetCompareTypes() map[string][]Compare {
	fullCompare := []Compare{GreaterThan, GreaterEqual, LowerThan, LowerEqual, Equal, NotEqual, Range}
	equalCompare := []Compare{Equal, NotEqual}
	stringCompare := []Compare{Equal, NotEqual}

	return map[string][]Compare{
		Float:    fullCompare,
		Integer:  fullCompare,
		DateTime: fullCompare,
		String:   stringCompare,
		Boolean:  equalCompare,
		Time:     fullCompare,
		Date:     fullCompare,
	}
}
