package serializer

const (
	NoErr = iota
	ParamErr
	RedisErr
)

type Response struct {
	Code  int         `json:"code"`
	Data  interface{} `json:"data,omitempty"`
	Msg   string      `json:"msg"`
	Error string      `json:"error,omitempty",`
}

func BuildResponse(msg string) Response {
	return Response{
		Code: 0,
		Msg:  msg,
	}
}
func Err(errtype int, err error) Response {
	return Response{
		Code:  errtype,
		Error: err.Error(),
	}
}
