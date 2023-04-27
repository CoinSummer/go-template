package go_template

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"math/big"
	"strconv"
	"strings"
	"testing"
	"time"
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
	data, _ := json.Marshal(res)
	assert.Equal(t, `"1"`, string(data))
}

func TestStringDecimal(t *testing.T) {
	got, err := NewExprFragment(`"8912239900000000000" / 1e18`, NewOperatorsMgr(), NewFnMgr())
	if err != nil {
		t.Error(err)
	}
	res, err := got.Eval("", config)
	if err != nil {
		t.Error(err)
	}

	data, _ := json.Marshal(res)
	assert.Equal(t, `"8.91"`, string(data))
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
	data, _ := json.Marshal(res)
	assert.Equal(t, `"3"`, string(data))
}

func TestBigInt(t *testing.T) {
	got, err := NewExprFragment(`$value / 1e18`, NewOperatorsMgr(), NewFnMgr())
	if err != nil {
		t.Error(err)
	}
	d, _ := big.NewInt(0).SetString("8912239900000000000", 10)
	env := map[string]interface{}{
		"value": d,
	}
	s, _ := json.Marshal(env)
	res, err := got.Eval(string(s), config)
	if err != nil {
		t.Error(err)
	}
	data, _ := json.Marshal(res)
	assert.Equal(t, `"8.91"`, string(data))
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
	data, _ := json.Marshal(res)
	assert.Equal(t, `"2"`, string(data))
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
	data, _ := json.Marshal(res)
	assert.Equal(t, `"3"`, string(data))
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
	data, _ := json.Marshal(res)
	assert.Equal(t, `"3"`, string(data))
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
	data, _ := json.Marshal(res)
	assert.Equal(t, `"5"`, string(data))
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
	data, _ := json.Marshal(res)
	assert.Equal(t, `"3"`, string(data))
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
	data, _ := json.Marshal(res)
	assert.Equal(t, `"3"`, string(data))
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
	data, _ := json.Marshal(res)
	assert.Equal(t, `"1"`, string(data))
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
	data, _ := json.Marshal(res)
	assert.Equal(t, `"1.1"`, string(data))
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
func TestThousandSep(t *testing.T) {
	got, err := NewExprFragment(`1000000000`, NewOperatorsMgr(), NewFnMgr())
	if err != nil {
		t.Fatal(err)
	}
	v, err := got.Eval(``, config)
	s, _ := json.Marshal(v)
	assert.Equal(t, string(s), `"1,000,000,000"`)
}
func TestThousandSepFloat(t *testing.T) {
	got, err := NewExprFragment(`1000000000.12332100000`, NewOperatorsMgr(), NewFnMgr())
	if err != nil {
		t.Fatal(err)
	}
	v, err := got.Eval(``, config)
	s, _ := json.Marshal(v)
	assert.Equal(t, string(s), `"1,000,000,000.12"`)
}
func TestThousandSepFloat3(t *testing.T) {
	got, err := NewExprFragment(`100.12332100000`, NewOperatorsMgr(), NewFnMgr())
	if err != nil {
		t.Fatal(err)
	}
	v, err := got.Eval(``, config)
	s, _ := json.Marshal(v)
	assert.Equal(t, string(s), `"100.12"`)
}
func TestThousandSepFloatLt1(t *testing.T) {
	got, err := NewExprFragment(`0.0000012332100000`, NewOperatorsMgr(), NewFnMgr())
	if err != nil {
		t.Fatal(err)
	}
	v, err := got.Eval(``, config)
	s, _ := json.Marshal(v)
	assert.Equal(t, `"0.0000012"`, string(s))
}
func TestThousandSepFloate18(t *testing.T) {
	got, err := NewExprFragment(`100000/1e4`, NewOperatorsMgr(), NewFnMgr())
	if err != nil {
		t.Fatal(err)
	}
	v, err := got.Eval(``, config)
	s, _ := json.Marshal(v)
	assert.Equal(t, `"10"`, string(s))
}
