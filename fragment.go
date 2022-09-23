package go_template

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/dop251/goja/ast"
	astParser "github.com/dop251/goja/parser"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
)

type IFragment interface {
	Eval() (interface{}, error)
	RawContent() string
}

type PlainFragment struct {
	Content string
}

func NewPlainFragment(text string) *PlainFragment {
	return &PlainFragment{
		Content: text,
	}
}

func (p *PlainFragment) Eval() (interface{}, error) {
	result := p.Content
	result = strings.ReplaceAll(result, "{{", "{")
	result = strings.ReplaceAll(result, "}}", "}")
	return result, nil
}

func (p *PlainFragment) RawContent() string {
	return p.Content
}

// -------------------------------------------------------------

type ExprFragment struct {
	Content string
	Ctx     string
	OpMgr   *OperatorsMgr
	FnMgr   *FnMgr
}

func NewExprFragment(text string, ctx string, opMgr *OperatorsMgr, fnMgr *FnMgr) *ExprFragment {
	return &ExprFragment{
		Content: text,
		Ctx:     ctx,
		OpMgr:   opMgr,
		FnMgr:   fnMgr,
	}
}

func (p *ExprFragment) RawContent() string {
	return p.Content
}
func ErrFMsg(format string, a ...any) error {
	logrus.Errorf(format, a...)
	return fmt.Errorf(format, a...)
}

func (f *ExprFragment) EvalBin(arg1, arg2 ast.Expression, op string) (interface{}, error) {
	arg1Value, err := f.EvalExpr(arg1)
	if err != nil {
		return arg1Value, err
	}
	arg2Value, err := f.EvalExpr(arg2)
	if err != nil {
		return arg2Value, err
	}
	operator := f.OpMgr.GetFunc(op)
	if operator == nil {
		return nil, ErrFMsg("operator not found: %s", op)
	}
	result, err := operator(arg1Value, arg2Value)
	if err != nil {
		return nil, ErrFMsg("operator error: %s", err)
	}

	return result, nil
}

func (f *ExprFragment) EvalCall(funcName string, args []ast.Expression) (interface{}, error) {
	fn := f.FnMgr.GetFunc(funcName)
	if fn == nil {
		return nil, ErrFMsg("func not found: %s", funcName)
	}
	var argsValue []interface{}
	for _, arg := range args {
		argValue, err := f.EvalExpr(arg)
		if err != nil {
			return nil, ErrFMsg("failed eval function args: %s, err: %s", funcName, err)
		}
		argsValue = append(argsValue, argValue)
	}
	result, err := fn(argsValue)
	if err != nil {
		return nil, ErrFMsg("failed eval function: %s, err: %s", funcName, err)
	}
	return result, nil
}

func (f *ExprFragment) Decimalize(value interface{}) interface{} {
	switch tp := value.(type) {
	case float32:
		return decimal.NewFromFloat32(tp)
	case float64:
		return decimal.NewFromFloat(tp)
	case int32:
		return decimal.NewFromInt32(tp)
	case int64:
		return decimal.NewFromInt(tp)
	default:
		return value
	}

}
func (f *ExprFragment) EvalExpr(expr ast.Expression) (interface{}, error) {
	switch expr := expr.(type) {
	case *ast.Identifier:
		name := expr.Name.String()
		if strings.HasPrefix(name, "$") {
			name = strings.TrimPrefix(name, "$")
		} else {
			return nil, ErrFMsg("unknown variable: %s", name)
		}
		value := gjson.Get(f.Ctx, name).Value()
		return f.Decimalize(value), nil
	case *ast.DotExpression:
		leftValue, err := f.EvalExpr(expr.Left)
		if err != nil {
			return nil, err
		}
		jStr, err := json.Marshal(leftValue)
		if err != nil {
			return nil, ErrFMsg("failed marshal dot left: %s err: %s", leftValue, err)
		}
		value := gjson.Get(string(jStr), expr.Identifier.Name.String()).Value()
		return f.Decimalize(value), nil
	case *ast.BinaryExpression:
		return f.EvalBin(expr.Left, expr.Right, expr.Operator.String())
	case *ast.CallExpression:
		funcName, ok := expr.Callee.(*ast.Identifier)
		if !ok {
			return "", ErrFMsg("<function not found: %s>", funcName.Name.String())
		}
		return f.EvalCall(funcName.Name.String(), expr.ArgumentList)
	case *ast.NumberLiteral:
		switch expr.Value.(type) {
		case float32:
			return decimal.NewFromFloat32(expr.Value.(float32)), nil
		case float64:
			return decimal.NewFromFloat(expr.Value.(float64)), nil
		case int64:
			return decimal.NewFromInt(expr.Value.(int64)), nil
		case int32:
			return decimal.NewFromInt32(expr.Value.(int32)), nil
		default:
			logrus.Warnf("unknown number literal: %v", expr.Value)
			return decimal.NewFromInt(0), nil
		}
	case *ast.StringLiteral:
		return expr.Value.String(), nil
	case *ast.BooleanLiteral:
		return expr.Value, nil
	default:
		return nil, ErrFMsg("expr not supported: %s", reflect.TypeOf(expr).String())
		// 三目运算符
		// case *ast.ConditionalExpression:
		// 	// todo tri operator
		// 	_ = expr.(*ast.CallExpression)
	}
}

func (f *ExprFragment) EvalContent(content string) (interface{}, error) {
	// todo eval expr
	p, err := astParser.ParseFile(nil, "", content, 0)
	if err != nil {
		return content, ErrFMsg("failed parse expr: %s, err: %s", content, err)
	} else {
		if len(p.Body) != 1 {
			return content, ErrFMsg(" < only support ONE expr: %s > ", content)
		}
		bodyType := reflect.TypeOf(p.Body[0]).String()
		if bodyType != "*ast.ExpressionStatement" {
			return content, ErrFMsg("expr not support: %s", content)
		}
		result, err := f.EvalExpr(p.Body[0].(*ast.ExpressionStatement).Expression)

		if err != nil {
			return content, ErrFMsg("< eval expr err: %s, err: %s > ", content, err)
		} else {
			return result, nil
		}
	}
}
func (f *ExprFragment) Eval() (interface{}, error) {
	return f.EvalContent(f.Content)
}
