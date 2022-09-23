package go_template

import (
	"reflect"

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
					return nil, ErrFMsg("/ with NaN: %v", a)
				}
				b, ok := arg2.(decimal.Decimal)
				if !ok {
					return nil, ErrFMsg("/ to NaN: %v", b)
				}
				return a.Div(b), nil
			},
			"+": func(arg1, arg2 interface{}) (interface{}, error) {
				a, ok := arg1.(decimal.Decimal)
				if !ok {
					return nil, ErrFMsg("+ with NaN: %s", reflect.TypeOf(arg1).String())
				}
				b, ok := arg2.(decimal.Decimal)
				if !ok {
					return nil, ErrFMsg("+ to NaN: %v", reflect.TypeOf(arg2).String())
				}
				return a.Add(b), nil
			},
			"-": func(arg1, arg2 interface{}) (interface{}, error) {
				a, ok := arg1.(decimal.Decimal)
				if !ok {
					return nil, ErrFMsg("- with NaN: %v", a)
				}
				b, ok := arg2.(decimal.Decimal)
				if !ok {
					return nil, ErrFMsg("- to NaN: %v", b)
				}
				return a.Sub(b), nil
			},
			"*": func(arg1, arg2 interface{}) (interface{}, error) {
				a, ok := arg1.(decimal.Decimal)
				if !ok {
					return nil, ErrFMsg("* with NaN: %v", a)
				}
				b, ok := arg2.(decimal.Decimal)
				if !ok {
					return nil, ErrFMsg("* to NaN: %v", b)
				}
				return a.Div(b), nil
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
