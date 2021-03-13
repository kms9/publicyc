package ogin

type ErrMsg struct {
	Err   string      `json:"err,omitempty"`
	Msg   string      `json:"msg,omitempty"`
	Debug string      `json:"debug,omitempty"`
	Data  interface{} `json:"data,omitempty"`
}
