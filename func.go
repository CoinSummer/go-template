package go_template

import (
	"fmt"
	"time"

	"github.com/shopspring/decimal"
)

type IFn func(config *TemplateConfig, args []interface{}) (interface{}, error)

type FnMgr struct {
	Funcs map[string]IFn
}

func tryParseTime(timeStr string, config *TemplateConfig) (time.Time, error) {
	formats := []string{time.RFC3339Nano, time.RFC3339, time.RFC822, time.RFC822Z, time.RFC850, time.RFC1123Z, time.RFC1123, time.ANSIC, time.UnixDate, time.Layout, time.RubyDate}
	formats = append(formats, config.TimeFormat)
	for _, format := range formats {
		t, err := time.Parse(format, timeStr)
		if err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("unknown time format: %s", timeStr)
}

func withTimezone(config *TemplateConfig, args []interface{}) (interface{}, error) {
	var timeOffset = config.TimeOffset
	var timeFormat = config.TimeFormat
	if len(args) >= 2 {
		offset, ok := args[1].(decimal.Decimal)
		if !ok {
			return nil, ErrFMsg("format arg1 must be timezone: %v", args[1])
		}
		timeOffset = int(offset.IntPart())
	}
	if len(args) == 3 {
		format, ok := args[2].(string)
		if !ok {
			return nil, ErrFMsg("format arg2 must be format string: %v", args[2])
		}
		timeFormat = format
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
		dt, err = tryParseTime(dtStr, config)
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

	return FormatTime(dt, timeOffset, timeFormat), nil
}

func NewFnMgr() *FnMgr {
	return &FnMgr{
		Funcs: map[string]IFn{
			"round": func(config *TemplateConfig, args []interface{}) (interface{}, error) {
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
			"timezone":   withTimezone,
			"formatTime": withTimezone,
		},
	}
}

func (f *FnMgr) RegisterFunc(name string, fn IFn) {
	f.Funcs[name] = fn
}

func (f *FnMgr) GetFunc(name string) IFn {
	return f.Funcs[name]
}
