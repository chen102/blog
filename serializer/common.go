package serializer

const (
	NoErr = iota
	ParamErr
	RedisErr
	StrconvErr
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
	resp := Response{
		Code:  errtype,
		Error: err.Error(),
	}
	switch errtype {
	case 1:
		resp.Msg = "参数错误"

	case 2:
		resp.Msg = "Redis操作失败"
	case 3:
		resp.Msg = "类型转换错误"
	default:
		resp.Msg = "未知错误" //这确定不是坑自己？

	}
	return resp
}
