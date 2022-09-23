package go_template

import "github.com/shopspring/decimal"

type IFn func(args []interface{}) (interface{}, error)

type FnMgr struct {
	Funcs map[string]IFn
}

func NewFnMgr() *FnMgr {
	return &FnMgr{
		Funcs: map[string]IFn{
			"round": func(args []interface{}) (interface{}, error) {
				if len(args) != 2 {
					return nil, ErrFMsg("round only accept 2 arg, got: %d", len(args))
				}
				n, ok := args[0].(decimal.Decimal)
				if !ok {
					return nil, ErrFMsg("round with NaN: %v", n)
				}
				place, ok := args[1].(decimal.Decimal)
				if !ok {
					return nil, ErrFMsg("round to NaN: %v", place)
				}
				return n.Round(int32(place.BigInt().Int64())), nil
			},
		},
	}
}

func (f *FnMgr) RegisterFunc(name string, fn IFn) {
	f.Funcs[name] = fn
}

func (f *FnMgr) GetFunc(name string) IFn {
	return f.Funcs[name]
}
