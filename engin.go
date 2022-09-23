package go_template

type TemplateEngine struct {
	FnMgr        *FnMgr
	OperatorsMgr *OperatorsMgr
}

func NewTemplateEngine() *TemplateEngine {
	fm := NewFnMgr()
	om := NewOperatorsMgr()
	return &TemplateEngine{
		FnMgr:        fm,
		OperatorsMgr: om,
	}
}
