package go_template

import (
	"testing"

	"github.com/sirupsen/logrus"
)

func TestTemplate_Render(t *testing.T) {
	tp, _ := NewTemplate("xxx { round(1000000000001000000 / 10e18, 1)} yyy", nil)
	res, err := tp.Render(`{}`)
	if err != nil {
		t.Error(err)
	}
	logrus.Info(res)

}

func TestTemplate_Render_Variable(t *testing.T) {
	tp, _ := NewTemplate("{$a.b + $a.b - 1} xxx", nil)
	res, err := tp.Render(`{"a": {"b": 333}}`)
	if err != nil {
		t.Error(err)
	}
	logrus.Info(res)

}
func TestTemplate_Render_Invalid_Variable(t *testing.T) {
	tp, _ := NewTemplate("{$a.c} xxx", nil)
	res, err := tp.Render(`{"a": {"b": 333}}`)
	if err != nil {
		t.Error(err)
	}
	if res != "{$a.c} xxx" {
		t.Errorf("expect %s, got %s", "{$a.c} xxx", res)
	}
}

func TestSyntaxError(t *testing.T) {
	_, err := NewTemplate("{111a.c} xxx", nil)
	if err == nil {
		t.Error(err)
	}
}
