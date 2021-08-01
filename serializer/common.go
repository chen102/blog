package serializer

type Response struct {
	Code  int         `json:"code"`
	Data  interface{} `json:"data,omitempty"`
	Msg   string      `json:"msg"`
	Error string      `json:"error",omitempty`
}

func Err(err error) Response {
	return Response{
		Error: err.Error(),
	}
}
