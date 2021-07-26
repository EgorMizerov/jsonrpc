package jsonrpc

type rpcError int

const (
	ParseError          rpcError = -32700
	InvalidReqError     rpcError = -32600
	MethodNotFoundError rpcError = -32601
	InvalidParamsError  rpcError = -32602
	InternalErrorError  rpcError = -32603
	ServerError         rpcError = -32000
)

func (e rpcError) Message() (m string) {
	switch e {
	case -32700:
		m = "Parse error"
	case -32600:
		m = "Invalid Request"
	case -32601:
		m = "Method not found"
	case -32602:
		m = "Invalid params"
	case -32603:
		m = "Internal error"
	case -32000:
		m = "Server error"
	}

	return
}
