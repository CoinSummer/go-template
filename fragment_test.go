package go_template

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/shopspring/decimal"
)

func TestSimpleExprFragment(t *testing.T) {
	got, err := NewExprFragment("1 + 1", NewOperatorsMgr(), NewFnMgr())
	if err != nil {
		t.Error(err)
	}
	res, err := got.Eval("")
	if err != nil {
		t.Error(err)
	}
	if !res.(decimal.Decimal).Equal(decimal.NewFromInt(2)) {
		t.Errorf("expect %v, got %v", 2, res.(decimal.Decimal))
	}
}

func TestMultipleOperatorExprFragment(t *testing.T) {
	got, err := NewExprFragment("1 + 1 + 1", NewOperatorsMgr(), NewFnMgr())
	if err != nil {
		t.Error(err)
	}
	res, err := got.Eval("")
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
	res, err := got.Eval("")
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
	res, err := got.Eval(`{"a": 3}`)
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
	res, err := got.Eval(`{"a": {"b": {"c":3}}}`)
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
	res, err := got.Eval(`{"a": {"b": [1,2]}}`)
	if err != nil {
		t.Fatal(err)
	}
	if !res.(decimal.Decimal).Equal(decimal.NewFromInt(1)) {
		t.Errorf("expect %v, got %v", 1, res.(decimal.Decimal))
	}
}

func TestSimpleFuncFragment(t *testing.T) {
	got, err := NewExprFragment("round(1.12,1.1)", NewOperatorsMgr(), NewFnMgr())
	if err != nil {
		t.Fatal(err)
	}
	res, err := got.Eval(``)
	if err != nil {
		t.Fatal(err)
	}
	if !res.(decimal.Decimal).Equal(decimal.NewFromFloat(1.1)) {
		t.Errorf("expect %v, got %v", 1.1, res.(decimal.Decimal))
	}
}

func TestTimezoneFuncFragment(t *testing.T) {
	now := time.Now()
	expectStr := now.In(time.FixedZone("x", 8*3600)).Format(time.RFC3339Nano)

	env1 := map[string]time.Time{
		"time": now,
	}
	env1Str, _ := json.Marshal(env1)

	got, err := NewExprFragment("timezone($time,8)", NewOperatorsMgr(), NewFnMgr())
	if err != nil {
		t.Fatal(err)
	}

	res, err := got.Eval(string(env1Str))
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
	res, err = got.Eval(string(env2Str))
	if err != nil {
		t.Fatal(err)
	}
	if !(res.(string) == expectStr) {
		t.Errorf("expect %s, got %s", expectStr, res)
	}
}
