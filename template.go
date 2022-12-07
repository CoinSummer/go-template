package go_template

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
)

type TemplateConfig struct {
	TimeOffset int
	TimeFormat string
}

type Template struct {
	templateText   string // plain text with variables embeded
	ctx            string // json string
	engine         *TemplateEngine
	parsedTemplate []IFragment
	TemplateConfig *TemplateConfig
}

func NewTemplate(text string, engine *TemplateEngine) (*Template, error) {
	return NewTemplateWithConfig(text, engine, nil)

}
func NewTemplateWithConfig(text string, engine *TemplateEngine, config *TemplateConfig) (*Template, error) {
	if engine == nil {
		engine = NewTemplateEngine()
	}
	if config == nil {
		config = &TemplateConfig{
			TimeOffset: 0,
			TimeFormat: strings.ReplaceAll(time.RFC3339, "T", " "),
		}
	}

	t := &Template{
		templateText:   text,
		engine:         engine,
		ctx:            "",
		TemplateConfig: config,
	}
	// parse template to fragments
	fragments, err := t.ParseFragments()
	if err != nil {
		return nil, err
	}
	t.parsedTemplate = fragments
	return t, nil
}

// can't unread more than once, use preifx to represent chars to unread
func (t *Template) ParsePlain(reader *strings.Reader, prefix string) (IFragment, error) {
	// assume start with plain
	text := strings.Builder{}
	text.WriteString(prefix)
	var ch rune
	var err error
	for {
		ch, _, err = reader.ReadRune()
		// EOF
		if err != nil {
			return NewPlainFragment(text.String()), nil
		}
		// expr next
		if ch == '{' {
			_ = reader.UnreadRune()
			return NewPlainFragment(text.String()), nil
		} else {
			text.WriteRune(ch)
		}
	}
}

// maybe a expr
func (t *Template) ParseMaybeExpr(reader *strings.Reader, prefix string) (IFragment, error) {
	// assume start with plain
	text := strings.Builder{}
	text.WriteString(prefix)

	var ch rune
	var err error
	// drop {
	_, _, _ = reader.ReadRune()

	bracketCount := 1

	for {
		ch, _, err = reader.ReadRune()
		// EOF
		if err != nil {
			// EOF before close bracket
			return NewPlainFragment(text.String()), nil
		}
		if ch == '{' {
			bracketCount += 1
		}
		if ch == '}' {
			bracketCount -= 1
		}

		if bracketCount == 0 {
			return NewExprFragment(text.String(), t.engine.OperatorsMgr, t.engine.FnMgr)
		} else {
			text.WriteRune(ch)
		}
	}
}

// split template to plain or expr part
func (t *Template) ParseFragments() ([]IFragment, error) {
	reader := strings.NewReader(t.templateText)
	// loop to read fragment， expr and plain Alternating
	fragments := []IFragment{}
	for {
		ch, _, err := reader.ReadRune()
		// EOF
		if err != nil {
			break
		}
		_ = reader.UnreadRune()
		if ch == '{' {
			f, err := t.ParseMaybeExpr(reader, "")
			if err != nil {
				return fragments, err
			}
			if f != nil {
				fragments = append(fragments, f)
			}
		} else {
			f, err := t.ParsePlain(reader, "")
			if err != nil {
				return fragments, err
			}
			if f != nil {
				fragments = append(fragments, f)
			}
		}
	}
	return fragments, nil

}

func (t *Template) RenderWithConfig(env string, config *TemplateConfig) (string, error) {
	if config == nil {
		config = t.TemplateConfig
	}
	t.ctx = env

	result := ""
	// eval fragments to string
	for _, f := range t.parsedTemplate {
		res, err := f.Eval(t.ctx, config)
		if err != nil {
			logrus.Warnf("failed eval template expression: %s", f.RawContent())
			//result += fmt.Sprintf("** %s ** ", err)
			result += "{" + f.RawContent() + "}"
			continue
		}
		if res == nil {
			result += "{" + f.RawContent() + "}"
			continue
		}

		j, err := json.Marshal(res)
		if err != nil {
			logrus.Warnf("failed marshal expr result: %s, err: %s", j, err)
			result += fmt.Sprintf("** %s ** ", err)
			continue
		}
		// concat fragments
		result += gjson.Parse(string(j)).String()
	}

	return result, nil
}
func (t *Template) Render(env string) (string, error) {
	return t.RenderWithConfig(env, t.TemplateConfig)
}
