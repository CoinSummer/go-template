package go_template

import (
	"time"
)

var Timezones = map[int]string{
	0:   "GMT",
	1:   "CET",
	2:   "EET",
	3:   "+03",
	4:   "+04",
	5:   "PKT",
	6:   "+06",
	7:   "+07",
	8:   "CST",
	9:   "JST",
	10:  "AEDT",
	11:  "+11",
	12:  "+12",
	13:  "+13",
	14:  "+14",
	-12: "-12",
	-11: "SST",
	-10: "HST",
	-9:  "AKST",
	-8:  "PST",
	-7:  "MST",
	-6:  "CST",
	-5:  "EST",
	-4:  "AST",
	-3:  "-3",
	-2:  "-2",
	-1:  "-1",
}

func FormatTime(t time.Time, offset int, format string) string {
	timezoneName, ok := Timezones[offset]
	if !ok {
		timezoneName = ""
	}

	return t.In(time.FixedZone(timezoneName, int(offset*3600))).Format(format)
}
