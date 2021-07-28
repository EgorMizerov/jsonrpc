package jsonrpc

import (
	"fmt"
	"github.com/valyala/fastjson"
	"io"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

// Server is a json-rpc server
type (
	Handler struct {
		handlers map[string]func(ctx *Context)
	}

	Method struct {
		fn func(ctx *Context)
	}

	Response struct {
		error  RpcError
		result *result
		id     interface{}
	}

	H map[string]interface{}
)

var (
	e1                 = ParseError
	parseErrorResponse = Response{
		error: &e1,
	}
	e2                  = ServerError
	serverErrorResponse = Response{
		error: &e2,
	}
)

func NewHandler() *Handler {
	return &Handler{
		handlers: make(map[string]func(ctx *Context)),
	}
}

// SetMethod is a function for setting a handler for a method
func (s *Handler) SetMethod(name string, fn func(ctx *Context)) {
	s.handlers[name] = fn
}

func (r *Response) String() string {
	str := "{\"jsonrpc\": \"2.0\", "
	if r.error != nil {
		str += "\"error\": {\"code\": " + strconv.Itoa(r.error.GetCode()) + ", \"message\": \"" + r.error.GetMessage() + "\"}"

		if r.id != nil {
			id, ok := r.id.(string)
			if ok {
				str += ", \"id\": \"" + id + "\"}"
				return str
			}
			id = fmt.Sprintf("%v", r.id)
			str += ", \"id\": " + id + "}"
			return str
		}
		str += ", \"id\": " + "null" + "}"
	}

	if r.result != nil {
		str += "\"result\": " + r.result.String()
		if r.id != nil {
			id, ok := r.id.(string)
			if ok {
				str += ", \"id\": \"" + id + "\"}"
				return str
			}
			id = fmt.Sprintf("%d", r.id)
			str += ", \"id\": " + id + "}"
			return str
		}
		str += ", \"id\": " + "null" + "}"
	}

	return str
}

func joinResponses(responses []*Response) (out string) {
	out = "["
	for i, response := range responses {
		if i == 0 {
			out += response.String()
			continue
		}
		out += ", " + response.String()
	}
	out += "]"
	return
}

func (s *Handler) WebSocket(r *http.Request, w http.ResponseWriter, in string) (string, error) {
	var p fastjson.Parser
	v, err := p.ParseBytes([]byte(in))
	if err != nil {
		if strings.Contains(err.Error(), "cannot parse JSON") {
			return parseErrorResponse.String(), nil
		}
		return "", err
	}

	if v.Type() == fastjson.TypeArray {
		var wg sync.WaitGroup
		values, _ := v.Array()
		var count int
		var responses []*Response

		wg.Add(len(values))
		for _, value := range values {
			var response Response
			responses = append(responses, &response)
			go s.handler(r, w, value, &response, &wg)
			count++
		}
		wg.Wait()

		out := joinResponses(responses)
		return out, nil
	}

	if v.Type() == fastjson.TypeObject {
		var response Response
		s.handler(r, w, v, &response, nil)
		return response.String(), nil
	}

	return "", err
}

func (s *Handler) RPC(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var p fastjson.Parser

	w.Header().Add("Content-Type", "application/json")

	data, err := io.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte(serverErrorResponse.String()))
		return
	}

	v, err := p.ParseBytes(data)
	if err != nil {
		if strings.Contains(err.Error(), "cannot parse JSON") {
			w.Write([]byte(parseErrorResponse.String()))
			return
		}
		return
	}

	if v.Type() == fastjson.TypeArray {
		var wg sync.WaitGroup
		values, _ := v.Array()
		var count int
		var responses []*Response

		wg.Add(len(values))
		for _, value := range values {
			var response Response
			responses = append(responses, &response)
			go s.handler(r, w, value, &response, &wg)
			count++
		}
		wg.Wait()

		out := joinResponses(responses)
		w.Write([]byte(out))
	}

	if v.Type() == fastjson.TypeObject {
		var response Response
		s.handler(r, w, v, &response, nil)
		w.Write([]byte(response.String()))
	}
}

func (s *Handler) handler(r *http.Request, w http.ResponseWriter, v *fastjson.Value, out *Response, wg *sync.WaitGroup) {
	if wg != nil {
		defer wg.Done()
	}
	ctx := newContext(r, out)
	rpcv := v.Get("jsonrpc")
	if rpcv == nil || rpcv.String() != "\"2.0\"" {
		err := InvalidReqError
		out.error = &err
		return
	}

	method := v.Get("method")
	if method == nil {
		err := InvalidReqError
		out.error = &err
		return
	}
	if method.Type() != fastjson.TypeString {
		err := InvalidReqError
		out.error = &err
		return
	}

	reqID := v.Get("id")
	if reqID != nil {
		switch reqID.Type() {
		case fastjson.TypeNumber:
			id, _ := reqID.Int()
			out.id = id
		case fastjson.TypeString:
			id := strings.Split(reqID.String(), "\"")[1]
			out.id = id
		}
	}

	fn, ok := s.handlers[strings.Split(method.String(), "\"")[1]]
	if !ok {
		err := MethodNotFoundError
		out.error = &err
		return
	}
	ctx.Set("params", v.Get("params"))
	fn(ctx)
}
