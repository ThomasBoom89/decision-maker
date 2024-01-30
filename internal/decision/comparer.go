package decision

type Compare string

const (
	GreaterThan  Compare = "gt"
	GreaterEqual Compare = "ge"
	LowerThan    Compare = "lt"
	LowerEqual   Compare = "le"
	Equal        Compare = "eq"
	NotEqual     Compare = "ne"
)

const (
	Float    = "float"
	Integer  = "int"
	DateTime = "datetime"
	String   = "string"
	Boolean  = "bool"
)

type Comparer interface {
	CompareInt(schmalue int64, value int64, compare Compare) bool
	CompareFloat(schmalue float64, value float64, compare Compare) bool
	CompareString(schmalue string, value string, compare Compare) bool
	CompareBool(schmalue bool, value bool, compare Compare) bool
	CompareDateTime(schmalue int64, value int64, compare Compare) bool
}

type TestConfigurationComparer struct {
	int    FullCompare[int64]
	float  FullCompare[float64]
	string EqualCompare[string]
	bool   EqualCompare[bool]
}

func (T *TestConfigurationComparer) CompareInt(schmalue int64, value int64, compare Compare) bool {
	return T.int.CompareFull(schmalue, value, compare)

}

func (T *TestConfigurationComparer) CompareFloat(schmalue float64, value float64, compare Compare) bool {
	return T.float.CompareFull(schmalue, value, compare)
}

func (T *TestConfigurationComparer) CompareString(schmalue string, value string, compare Compare) bool {
	return T.string.CompareEqual(schmalue, value, compare)
}

func (T *TestConfigurationComparer) CompareBool(schmalue bool, value bool, compare Compare) bool {
	return T.bool.CompareEqual(schmalue, value, compare)
}

func (T *TestConfigurationComparer) CompareDateTime(schmalue int64, value int64, compare Compare) bool {
	return T.int.CompareFull(schmalue, value, compare)
}

type RealComparer struct {
	int    FullCompare[int64]
	float  FullCompare[float64]
	string EqualCompare[string]
	bool   EqualCompare[bool]
}

func (R *RealComparer) CompareInt(schmalue int64, value int64, compare Compare) bool {
	return R.int.CompareFull(schmalue, value, compare)
}

func (R *RealComparer) CompareFloat(schmalue float64, value float64, compare Compare) bool {
	return R.float.CompareFull(schmalue, value, compare)
}

func (R *RealComparer) CompareString(schmalue string, value string, compare Compare) bool {
	return R.string.CompareEqual(schmalue, value, compare)
}

func (R *RealComparer) CompareBool(schmalue bool, value bool, compare Compare) bool {
	return R.bool.CompareEqual(schmalue, value, compare)
}

func (R *RealComparer) CompareDateTime(schmalue int64, value int64, compare Compare) bool {
	return R.int.CompareFull(schmalue, value, compare)
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

func (C *TestConfigurationFullCompare[T]) CompareFull(schmalue T, value T, compare Compare) bool {
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
	fullCompare := []Compare{GreaterThan, GreaterEqual, LowerThan, LowerEqual, Equal, NotEqual}
	equalCompare := []Compare{Equal, NotEqual}

	return map[string][]Compare{
		Float:    fullCompare,
		Integer:  fullCompare,
		DateTime: fullCompare,
		String:   equalCompare,
		Boolean:  equalCompare,
	}
}
