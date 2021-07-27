package jsonrpc

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
)

type bindType int

const (
	BindJSONType bindType = iota
)

func bind(obj interface{}, bindType bindType, ctx *Context) error {
	params, err := ctx.Params()
	if err != nil {
		return err
	}
	switch bindType {
	case BindJSONType:
		return decodeJSON(params.MarshalTo(nil), obj)
	}
	return nil
}

func validate(obj interface{}) error {
	v := validator.New()
	return v.Struct(obj)
}

func decodeJSON(data []byte, obj interface{}) error {
	err := json.Unmarshal(data, obj)
	if err != nil {
		return err
	}

	return validate(obj)
}
