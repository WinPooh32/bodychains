// Code generated by gotemplate. DO NOT EDIT.

package requests

// template type Template()

type (
	GetStateByHash struct {
		ArgsGetStateByHash
		RespGetStateByHash
	}

	ArgsGetStateByHash struct {
	}

	RespGetStateByHash struct {
	}
)

func (t *GetStateByHash) Args() interface{} {
	var args interface{} = &t.ArgsGetStateByHash
	return args
}

func (t *GetStateByHash) Resp() interface{} {
	var resp interface{} = &t.RespGetStateByHash
	return resp
}
