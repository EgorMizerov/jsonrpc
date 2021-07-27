package jsonrpc

import (
	"fmt"
	"github.com/valyala/fastjson"
	"io"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Server is a json-rpc server
type (
	Server struct {
		server   *http.Server
		handlers map[string]func(ctx *Context)
	}

	Method struct {
		fn func(ctx *Context)
	}

	Response struct {
		error  *rpcError
		result *result
		id     interface{}
	}
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

// NewServer is a constructor for type Server
func NewServer(params ...ServerParam) (server *Server, err error) {
	server = &Server{
		server: &http.Server{
			Addr:         "localhost:8000",
			ReadTimeout:  time.Second * 10,
			WriteTimeout: time.Second * 10,
		},
		handlers: make(map[string]func(ctx *Context)),
	}
	server.server.Handler = server

	for _, param := range params {
		err = param.fn(server)
		if err != nil {
			return
		}
	}
	return
}

// Run is a function for listen and serve http server
func (s *Server) Run() error {
	return s.server.ListenAndServe()
}

// SetMethod is a function for setting a handler for a method
func (s *Server) SetMethod(name string, fn func(ctx *Context)) {
	s.handlers[name] = fn
}

func (r *Response) String() string {
	str := "{\"jsonrpc\": \"2.0\", "
	if r.error != nil {
		str += "\"error\": {\"code\": " + strconv.Itoa(int(*r.error)) + ", \"message\": \"" + r.error.Message() + "\"}"

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

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

func (s *Server) handler(r *http.Request, w http.ResponseWriter, v *fastjson.Value, out *Response, wg *sync.WaitGroup) {
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
