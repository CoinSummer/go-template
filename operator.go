package go_template

import (
	"github.com/shopspring/decimal"
)

type IOperator func(arg1, arg2 interface{}) (interface{}, error)

type OperatorsMgr struct {
	Operators map[string]IOperator
}

func decimalize(arg interface{}) (decimal.Decimal, error) {
	var a decimal.Decimal
	var ok bool
	var err error
	a, ok = arg.(decimal.Decimal)
	if !ok {
		if s, ok := arg.(string); ok {
			a, err = decimal.NewFromString(s)
			if err != nil {
				return decimal.Decimal{}, ErrFMsg("/ with NaN: %v", arg)
			} else {
				return a, nil
			}
		} else {
			return decimal.Decimal{}, ErrFMsg("/ with NaN: %v", arg)
		}
	} else {
		return a, nil
	}
}

func NewOperatorsMgr() *OperatorsMgr {
	return &OperatorsMgr{
		Operators: map[string]IOperator{
			"/": func(arg1, arg2 interface{}) (interface{}, error) {
				var a, b decimal.Decimal
				var err error
				a, err = decimalize(arg1)
				if err != nil {
					return nil, err
				}
				b, err = decimalize(arg2)
				if err != nil {
					return nil, err
				}
				return a.Div(b), nil
			},
			"+": func(arg1, arg2 interface{}) (interface{}, error) {
				var a, b decimal.Decimal
				var err error
				a, err = decimalize(arg1)
				if err != nil {
					return nil, err
				}
				b, err = decimalize(arg2)
				if err != nil {
					return nil, err
				}
				return a.Add(b), nil
			},
			"-": func(arg1, arg2 interface{}) (interface{}, error) {
				var a, b decimal.Decimal
				var err error
				a, err = decimalize(arg1)
				if err != nil {
					return nil, err
				}
				b, err = decimalize(arg2)
				if err != nil {
					return nil, err
				}
				return a.Sub(b), nil
			},
			"*": func(arg1, arg2 interface{}) (interface{}, error) {
				var a, b decimal.Decimal
				var err error
				a, err = decimalize(arg1)
				if err != nil {
					return nil, err
				}
				b, err = decimalize(arg2)
				if err != nil {
					return nil, err
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
