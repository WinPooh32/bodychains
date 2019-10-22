package template

// template type Template()

type (
	Template struct {
		ArgsTemplate
		RespTemplate
	}

	ArgsTemplate struct {
	}

	RespTemplate struct {
	}
)

func (t *Template) Args() interface{} {
	var args interface{} = &t.ArgsTemplate
	return args
}

func (t *Template) Resp() interface{} {
	var resp interface{} = &t.RespTemplate
	return resp
}
