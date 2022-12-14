package go_template

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"math/big"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/shopspring/decimal"
)

var config = &TemplateConfig{
	TimeOffset: 0,
	TimeFormat: strings.ReplaceAll(time.RFC3339, "T", " "),
}

func TestNumberLiteral(t *testing.T) {
	got, err := NewExprFragment(`1`, NewOperatorsMgr(), NewFnMgr())
	if err != nil {
		t.Error(err)
	}
	res, err := got.Eval("", config)
	if err != nil {
		t.Error(err)
	}
	if !res.(decimal.Decimal).Equal(decimal.NewFromInt(1)) {
		t.Errorf("expect %v, got %v", 1, res.(decimal.Decimal))
	}
}
func TestStringDecimal(t *testing.T) {
	got, err := NewExprFragment(`"8902239900000000000" / 1e18`, NewOperatorsMgr(), NewFnMgr())
	if err != nil {
		t.Error(err)
	}
	res, err := got.Eval("", config)
	if err != nil {
		t.Error(err)
	}
	if !res.(decimal.Decimal).Equal(decimal.NewFromFloat(8.9022399)) {
		t.Errorf("expect %v, got %v", 2, res.(decimal.Decimal))
	}
}

func TestInt(t *testing.T) {
	got, err := NewExprFragment(`9 / 3`, NewOperatorsMgr(), NewFnMgr())
	if err != nil {
		t.Error(err)
	}
	res, err := got.Eval(``, config)
	if err != nil {
		t.Error(err)
	}
	if !res.(decimal.Decimal).Equal(decimal.NewFromFloat(3)) {
		t.Errorf("expect %v, got %v", 3, res.(decimal.Decimal))
	}
}

func TestBigInt(t *testing.T) {
	got, err := NewExprFragment(`$value / 1e18`, NewOperatorsMgr(), NewFnMgr())
	if err != nil {
		t.Error(err)
	}
	d, _ := big.NewInt(0).SetString("8902239900000000000", 10)
	env := map[string]interface{}{
		"value": d,
	}
	s, _ := json.Marshal(env)
	res, err := got.Eval(string(s), config)
	if err != nil {
		t.Error(err)
	}
	if !res.(decimal.Decimal).Equal(decimal.NewFromFloat(8.9022399)) {
		t.Errorf("expect %v, got %v", 2, res.(decimal.Decimal))
	}
}
func TestSimpleExprFragment(t *testing.T) {
	got, err := NewExprFragment("1 + 1", NewOperatorsMgr(), NewFnMgr())
	if err != nil {
		t.Error(err)
	}
	res, err := got.Eval("", config)
	if err != nil {
		t.Error(err)
	}
	if !res.(decimal.Decimal).Equal(decimal.NewFromInt(2)) {
		t.Errorf("expect %v, got %v", 2, res.(decimal.Decimal))
	}
}

func TestAutoCastExprFragment(t *testing.T) {
	got, err := NewExprFragment(`"1" + "1" + 1`, NewOperatorsMgr(), NewFnMgr())
	if err != nil {
		t.Error(err)
	}
	res, err := got.Eval("", config)
	if err != nil {
		t.Error(err)
	}
	if !res.(decimal.Decimal).Equal(decimal.NewFromInt(3)) {
		t.Errorf("expect %v, got %v", 3, res.(decimal.Decimal))
	}
}
func TestMultipleOperatorExprFragment(t *testing.T) {
	got, err := NewExprFragment("1 + 1 + 1", NewOperatorsMgr(), NewFnMgr())
	if err != nil {
		t.Error(err)
	}
	res, err := got.Eval("", config)
	if err != nil {
		t.Error(err)
	}
	if !res.(decimal.Decimal).Equal(decimal.NewFromInt(3)) {
		t.Errorf("expect %v, got %v", 3, res.(decimal.Decimal))
	}
}

func TestParOperatorExprFragment(t *testing.T) {
	got, err := NewExprFragment("1 + (2 * 2)", NewOperatorsMgr(), NewFnMgr())
	if err != nil {
		t.Error(err)
	}
	res, err := got.Eval("", config)
	if err != nil {
		t.Error(err)
	}
	if !res.(decimal.Decimal).Equal(decimal.NewFromInt(5)) {
		t.Errorf("expect %v, got %v", 5, res.(decimal.Decimal))
	}
}

func TestSimpleVariableExprFragment(t *testing.T) {
	got, err := NewExprFragment("$a", NewOperatorsMgr(), NewFnMgr())
	if err != nil {
		t.Error(err)
	}
	res, err := got.Eval(`{"a": 3}`, config)
	if err != nil {
		t.Error(err)
	}
	if !res.(decimal.Decimal).Equal(decimal.NewFromInt(3)) {
		t.Errorf("expect %v, got %v", 3, res.(decimal.Decimal))
	}
}

func TestNestedVariableExprFragment(t *testing.T) {
	got, err := NewExprFragment("$a.b.c", NewOperatorsMgr(), NewFnMgr())
	if err != nil {
		t.Error(err)
	}
	res, err := got.Eval(`{"a": {"b": {"c":3}}}`, config)
	if err != nil {
		t.Error(err)
	}
	if !res.(decimal.Decimal).Equal(decimal.NewFromInt(3)) {
		t.Errorf("expect %v, got %v", 3, res.(decimal.Decimal))
	}
}

func TestArrayVariableExprFragment(t *testing.T) {
	got, err := NewExprFragment("$a.b[0]", NewOperatorsMgr(), NewFnMgr())
	if err != nil {
		t.Fatal(err)
	}
	res, err := got.Eval(`{"a": {"b": [1,2]}}`, config)
	if err != nil {
		t.Fatal(err)
	}
	if !res.(decimal.Decimal).Equal(decimal.NewFromInt(1)) {
		t.Errorf("expect %v, got %v", 1, res.(decimal.Decimal))
	}
}
func TestInExistVariableExprFragment(t *testing.T) {
	got, err := NewExprFragment("$a.c.d", NewOperatorsMgr(), NewFnMgr())
	if err != nil {
		t.Fatal(err)
	}
	_, err = got.Eval(`{"a": {"b": [1,2]}}`, config)
	if err == nil {
		t.Fatal(err)
	}
}
func TestBracket(t *testing.T) {
	got, err := NewExprFragment(`$a['b.c']`, NewOperatorsMgr(), NewFnMgr())
	if err != nil {
		t.Fatal(err)
	}
	v, err := got.Eval(`{"a": {"b.c": [1,2]}}`, config)
	s, _ := json.Marshal(v)
	assert.Equal(t, string(s), "[1,2]")

	v, err = got.Eval(`{"a": {"c": 123}}`, config)
	s, _ = json.Marshal(v)
	assert.Equal(t, string(s), "null")

}
func TestDotPathBracket(t *testing.T) {
	got, err := NewExprFragment(`$a['b/...c']`, NewOperatorsMgr(), NewFnMgr())
	if err != nil {
		t.Fatal(err)
	}
	v, err := got.Eval(`{"a": {"b/...c": [1,2]}}`, config)
	s, _ := json.Marshal(v)
	assert.Equal(t, string(s), "[1,2]")
}

func TestSimpleFuncFragment(t *testing.T) {
	got, err := NewExprFragment("round(1.12,1.1)", NewOperatorsMgr(), NewFnMgr())
	if err != nil {
		t.Fatal(err)
	}
	res, err := got.Eval(``, config)
	if err != nil {
		t.Fatal(err)
	}
	if !res.(decimal.Decimal).Equal(decimal.NewFromFloat(1.1)) {
		t.Errorf("expect %v, got %v", 1.1, res.(decimal.Decimal))
	}
}

func TestTimezoneFuncFragment(t *testing.T) {
	now := time.Now()
	expectStr := strings.ReplaceAll(now.In(time.FixedZone("x", 8*3600)).Format(time.RFC3339), "T", " ")

	env1 := map[string]interface{}{
		"time":   now,
		"ts":     now.UnixMilli(),
		"ts_str": strconv.FormatInt(now.UnixMilli(), 10),
	}
	env1Str, _ := json.Marshal(env1)
	got, err := NewExprFragment("timezone($ts,8)", NewOperatorsMgr(), NewFnMgr())
	if err != nil {
		t.Fatal(err)
	}
	r, err := got.Eval(string(env1Str), config)
	if r.(string) != expectStr {
		t.Errorf("expect %s, got: %s", expectStr, r)
	}

	got, err = NewExprFragment("timezone($ts_str,8)", NewOperatorsMgr(), NewFnMgr())
	if err != nil {
		t.Fatal(err)
	}
	r, err = got.Eval(string(env1Str), config)

	got, err = NewExprFragment("timezone($time,8)", NewOperatorsMgr(), NewFnMgr())
	if err != nil {
		t.Fatal(err)
	}

	res, err := got.Eval(string(env1Str), config)
	if err != nil {
		t.Fatal(err)
	}
	if !(res.(string) == expectStr) {
		t.Errorf("expect %s, got %s", expectStr, res)
	}

	env2 := map[string]string{
		"time": now.Format(time.RFC3339Nano),
	}
	env2Str, _ := json.Marshal(env2)
	res, err = got.Eval(string(env2Str), config)
	if err != nil {
		t.Fatal(err)
	}
	if !(res.(string) == expectStr) {
		t.Errorf("expect %s, got %s", expectStr, res)
	}
}
func TestTimezoneFuncCustomConfig(t *testing.T) {
	format := "2006-01-02 X 15:04:05Z07:00"
	config := &TemplateConfig{
		TimeOffset: 7,
		TimeFormat: format,
	}
	now := time.Now()
	expectStr := now.In(time.FixedZone("x", 8*3600)).Format(format)

	env1 := map[string]interface{}{
		"time": now,
	}

	env1Str, _ := json.Marshal(env1)

	got, err := NewExprFragment(`formatTime($time,8, "2006-01-02 X 15:04:05Z07:00")`, NewOperatorsMgr(), NewFnMgr())
	if err != nil {
		t.Fatal(err)
	}

	res, err := got.Eval(string(env1Str), config)
	if err != nil {
		t.Fatal(err)
	}

	if !(res.(string) == expectStr) {
		t.Errorf("expect %s, got %s", expectStr, res)
	}
}
