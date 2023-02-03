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
	Eval(ctx string, config *TemplateConfig) (interface{}, error)
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

func (p *PlainFragment) Eval(_ string, _ *TemplateConfig) (interface{}, error) {
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
	Content string // without {}
	Ast     *ast.Program
	Ctx     string
	OpMgr   *OperatorsMgr
	FnMgr   *FnMgr
}

func NewExprFragment(text string, opMgr *OperatorsMgr, fnMgr *FnMgr) (*ExprFragment, error) {
	f := &ExprFragment{
		Content: text,
		OpMgr:   opMgr,
		FnMgr:   fnMgr,
	}
	p, err := astParser.ParseFile(nil, "", text, 0)
	if err != nil {
		return nil, ErrFMsg("failed parse expr: %s, err: %s", text, err)
	}
	f.Ast = p
	if len(p.Body) != 1 {
		return nil, ErrFMsg(" < only support ONE expr: %s > ", text)
	}
	bodyType := reflect.TypeOf(f.Ast.Body[0]).String()
	if bodyType != "*ast.ExpressionStatement" {
		return nil, ErrFMsg("expr not support: %s", text)
	}
	return f, nil
}

func (p *ExprFragment) RawContent() string {
	return p.Content
}
func ErrFMsg(format string, a ...interface{}) error {
	logrus.Errorf(format, a...)
	return fmt.Errorf(format, a...)
}

// 二元操作符
func (f *ExprFragment) EvalBin(arg1, arg2 ast.Expression, op string, config *TemplateConfig) (interface{}, error) {
	arg1Value, err := f.EvalExpr(arg1, config)
	if err != nil {
		return arg1Value, err
	}
	arg2Value, err := f.EvalExpr(arg2, config)
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

func (f *ExprFragment) EvalCall(funcName string, args []ast.Expression, config *TemplateConfig) (interface{}, error) {
	fn := f.FnMgr.GetFunc(funcName)
	if fn == nil {
		return nil, ErrFMsg("func not found: %s", funcName)
	}
	var argsValue []interface{}
	for _, arg := range args {
		argValue, err := f.EvalExpr(arg, config)
		if err != nil {
			return nil, ErrFMsg("failed eval function args: %s, err: %s", funcName, err)
		}
		argsValue = append(argsValue, argValue)
	}
	result, err := fn(config, argsValue)
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
func (f *ExprFragment) EvalExpr(expr ast.Expression, config *TemplateConfig) (interface{}, error) {
	switch expr := expr.(type) {
	case *ast.Identifier:
		name := expr.Name.String()
		if strings.HasPrefix(name, "$") {
			name = strings.TrimPrefix(name, "$")
		} else {
			// 不支持变量
			return name, ErrFMsg("unsupported variable: %s", name)
		}
		value := gjson.Get(f.Ctx, name).Value()
		if value == nil {
			logrus.Warnf("variable %s not found in env %s: ", name, f.Ctx)
			return name, ErrFMsg("unknown variable: %s", name)
		}
		return f.Decimalize(value), nil

	case *ast.BracketExpression:
		leftValue, err := f.EvalExpr(expr.Left, config)
		if err != nil {
			return nil, err
		}
		jStr, err := json.Marshal(leftValue)
		if err != nil {
			return nil, ErrFMsg("failed marshal bracket left: %s err: %s", leftValue, err)
		}
		memberValue, err := f.EvalExpr(expr.Member, config)
		if err != nil {
			return nil, ErrFMsg("failed eval bracket member: %s err: %s", expr.Member, err)
		}
		var value interface{}
		text := string(jStr)
		switch m := memberValue.(type) {
		case decimal.Decimal:
			value = gjson.Get(text, m.BigInt().String()).Value()
		case string:
			// escaping path from gjson
			m = strings.ReplaceAll(m, ".", `\.`)
			value = gjson.Get(text, m).Value()
		default:
			return nil, ErrFMsg("index must be int or string, got: %s", reflect.TypeOf(m))
		}
		return f.Decimalize(value), nil
	case *ast.DotExpression:
		leftValue, err := f.EvalExpr(expr.Left, config)
		if err != nil {
			return nil, err
		}
		jStr, err := json.Marshal(leftValue)
		if err != nil {
			return nil, ErrFMsg("failed marshal dot left: %s err: %s", leftValue, err)
		}
		value := gjson.Get(string(jStr), expr.Identifier.Name.String()).Value()
		if value == nil {
			return nil, ErrFMsg("text %s not found in %s", expr.Identifier.Name.String(), string(jStr))
		}
		return f.Decimalize(value), nil
	case *ast.BinaryExpression:
		return f.EvalBin(expr.Left, expr.Right, expr.Operator.String(), config)
	case *ast.CallExpression:
		funcName, ok := expr.Callee.(*ast.Identifier)
		if !ok {
			return "", ErrFMsg("<function not found: %s>", funcName.Name.String())
		}
		return f.EvalCall(funcName.Name.String(), expr.ArgumentList, config)
	case *ast.NumberLiteral:
		d, err := decimal.NewFromString(fmt.Sprintf("%v", expr.Value))
		if err != nil {
			return nil, ErrFMsg("bad number literal %v", expr.Value)
		}
		return d, nil
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

func (f *ExprFragment) EvalContent(content string, config *TemplateConfig) (interface{}, error) {
	result, err := f.EvalExpr(f.Ast.Body[0].(*ast.ExpressionStatement).Expression, config)
	if err != nil {
		return content, ErrFMsg("< eval expr err: %s, err: %s > ", content, err)
	} else {
		return result, nil
	}
}
func (f *ExprFragment) Eval(ctx string, config *TemplateConfig) (interface{}, error) {
	f.Ctx = ctx
	result, err := f.EvalContent(f.Content, config)
	if err != nil {
		return result, err
	}

	switch r := result.(type) {
	case string:
		decRepr, err := decimal.NewFromString(r)
		if err != nil {
			return result, nil
		}
		return thousandSep(decRepr.String()), nil
	case decimal.Decimal:
		return thousandSep(r.String()), nil
	default:
		return result, nil
	}

}

func formatIntThousandSep(num string) string {
	var result []byte
	for i := 0; i < len(num); i++ {
		result = append([]byte{num[len(num)-i-1]}, result...)
		if (i+1)%3 == 0 {
			result = append([]byte{','}, result...)
		}
	}
	return string(result)
}

func thousandSep(num string) string {
	if strings.ContainsRune(num, '.') {
		items := strings.Split(num, ".")
		i := items[0]
		f := items[1]
		return formatIntThousandSep(i) + "." + f
	} else {
		return formatIntThousandSep(num)
	}
}
