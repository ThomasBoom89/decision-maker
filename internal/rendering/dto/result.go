package dto

import "github.com/ThomasBoom89/decision-maker/internal/decision"

type Result struct {
	ParameterName string
	TestValue     string
	ProductValue  string
	CompareType   string
	Comparer      decision.Compare
	Result        bool
}
