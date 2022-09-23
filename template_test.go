package go_template

import (
	"testing"

	"github.com/sirupsen/logrus"
)

func TestTemplate_Render(t *testing.T) {
	engine := NewTemplateEngine()
	tp := NewTemplate("xxx { round(1000000000001000000 / 10e18, 1)} yyy", "{}", engine)
	res, err := tp.Render()
	if err != nil {
		t.Error(err)
	}
	logrus.Info(res)

}

func TestTemplate_Render_Variable(t *testing.T) {
	engine := NewTemplateEngine()
	tp := NewTemplate("{$a.b + $a.b - 1} xxx", `{"a": {"b": 333}}`, engine)
	res, err := tp.Render()
	if err != nil {
		t.Error(err)
	}
	logrus.Info(res)

}
