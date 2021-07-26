package jsonrpc

import (
	"errors"
	"github.com/valyala/fastjson"
	"time"
)

type Context struct {
	values   map[interface{}]interface{}
	response *Response
}

func newContext(resp *Response) *Context {
	return &Context{
		values:   make(map[interface{}]interface{}),
		response: resp,
	}
}

func (c *Context) Deadline() (deadline time.Time, ok bool) {
	return time.Time{}, false
}

func (c *Context) Done() <-chan struct{} {
	return nil
}

func (c *Context) Err() error {
	return nil
}

func (c *Context) Set(key, value interface{}) {
	c.values[key] = value
}

func (c *Context) Value(key interface{}) interface{} {
	value, ok := c.values[key]
	if !ok {
		return nil
	}

	return value
}

func (c *Context) Params() (fastjson.Value, error) {
	val1 := c.Value("params")
	if val1 == nil {
		return fastjson.Value{}, errors.New("params not found")
	}
	val2, ok := val1.(*fastjson.Value)
	if !ok {
		return fastjson.Value{}, errors.New("failed to typecast")
	}

	return *val2, nil
}

func (c *Context) String(str string) {
	res := newStringResult(str)
	c.response.result = res
}

func (c *Context) Int(i int) {
	res := newIntResult(i)
	c.response.result = res
}

func (c *Context) Float(f float64) {
	res := newFloatResult(f)
	c.response.result = res
}

func (c *Context) Bool(b bool) {
	res := newBoolResult(b)
	c.response.result = res
}

func (c *Context) Array(arr interface{}) error {
	res, err := newArrayResult(arr)
	if err != nil {
		return err
	}
	c.response.result = res
	return nil
}

func (c *Context) Object(obj interface{}) error {
	res, err := newObjectResult(obj)
	if err != nil {
		return err
	}
	c.response.result = res
	return nil
}

func (c *Context) Error(err rpcError) {
	c.response.error = &err
}
