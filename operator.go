package go_template

import (
	"github.com/shopspring/decimal"
)

type IOperator func(arg1, arg2 interface{}) (interface{}, error)

type OperatorsMgr struct {
	Operators map[string]IOperator
}

func NewOperatorsMgr() *OperatorsMgr {
	return &OperatorsMgr{
		Operators: map[string]IOperator{
			"/": func(arg1, arg2 interface{}) (interface{}, error) {
				a, ok := arg1.(decimal.Decimal)
				if !ok {
					return nil, ErrFMsg("/ with NaN: %v", arg1)
				}
				b, ok := arg2.(decimal.Decimal)
				if !ok {
					return nil, ErrFMsg("/ to NaN: %v", arg2)
				}
				return a.Div(b), nil
			},
			"+": func(arg1, arg2 interface{}) (interface{}, error) {
				a, ok := arg1.(decimal.Decimal)
				if !ok {
					return nil, ErrFMsg("+ with NaN: %s", arg1)
				}
				b, ok := arg2.(decimal.Decimal)
				if !ok {
					return nil, ErrFMsg("+ to NaN: %v", arg1)
				}
				return a.Add(b), nil
			},
			"-": func(arg1, arg2 interface{}) (interface{}, error) {
				a, ok := arg1.(decimal.Decimal)
				if !ok {
					return nil, ErrFMsg("- with NaN: %v", arg1)
				}
				b, ok := arg2.(decimal.Decimal)
				if !ok {
					return nil, ErrFMsg("- to NaN: %v", arg2)
				}
				return a.Sub(b), nil
			},
			"*": func(arg1, arg2 interface{}) (interface{}, error) {
				a, ok := arg1.(decimal.Decimal)
				if !ok {
					return nil, ErrFMsg("* with NaN: %v", arg1)
				}
				b, ok := arg2.(decimal.Decimal)
				if !ok {
					return nil, ErrFMsg("* to NaN: %v", arg2)
				}
				return a.Mul(b), nil
			},
		},
	}
}

func (f *OperatorsMgr) RegisterFunc(name string, fn IOperator) {
	f.Operators[name] = fn
}

func (f *OperatorsMgr) GetFunc(name string) IOperator {
	return f.Operators[name]
}
