package go_template

import (
	"time"

	"github.com/shopspring/decimal"
)

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
			"timezone": func(args []interface{}) (interface{}, error) {
				if len(args) != 2 {
					return nil, ErrFMsg("timezone only accept 2 arg, got: %d", len(args))
				}
				var dt time.Time
				var err error
				dtStr, ok := args[0].(string)
				if !ok {
					// maybe timestamp int64
					if ts, ok := args[0].(decimal.Decimal); ok {
						if ts.GreaterThan(decimal.NewFromInt(10000000000)) {
							dt = time.UnixMilli(ts.IntPart())
						} else {
							dt = time.Unix(ts.IntPart(), 0)
						}
					} else {
						return nil, ErrFMsg("timezone arg0 must be string|int64|uint64: %v", dtStr)
					}
				} else {
					dt, err = time.Parse(time.RFC3339Nano, dtStr)
					if err != nil {
						ts, err := decimal.NewFromString(dtStr)
						if err != nil {
							return nil, ErrFMsg("timezone arg0 must be rfc3339nano format or timestamp: %s", dtStr)
						}
						if ts.GreaterThan(decimal.NewFromInt(10000000000)) {
							dt = time.UnixMilli(ts.IntPart())
						} else {
							dt = time.Unix(ts.IntPart(), 0)
						}
					}
				}
				offset, ok := args[1].(decimal.Decimal)
				if !ok {
					return nil, ErrFMsg("timezone arg1 must be number: %v", offset)
				}
				zonedDt := dt.In(time.FixedZone("custom", int(offset.BigInt().Int64())*3600))
				return zonedDt.Format(time.RFC3339Nano), nil
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
