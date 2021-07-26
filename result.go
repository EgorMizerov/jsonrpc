package jsonrpc

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

type (
	result struct {
		value      string
		typeResult typeResult
	}

	typeResult int
)

const (
	IntResultType typeResult = iota
	FloatResultType
	StringResultType
	BoolResultType
	ArrayResultType
	ObjectResultType
)

func (r *result) Type() typeResult {
	return r.typeResult
}
func (r *result) String() (out string) {
	switch r.typeResult {
	case StringResultType:
		out = "\"" + r.value + "\""
	case ArrayResultType & ObjectResultType & IntResultType & FloatResultType:
		out = r.value
	case BoolResultType:
		out = strings.Title(r.value)
	}
	return
}

func newStringResult(str string) *result {
	return &result{
		value:      str,
		typeResult: StringResultType,
	}
}

func newIntResult(dig int) *result {
	return &result{
		value:      strconv.Itoa(dig),
		typeResult: IntResultType,
	}
}

func newFloatResult(f float64) *result {
	return &result{
		value:      fmt.Sprintf("%f", f),
		typeResult: FloatResultType,
	}
}

func newBoolResult(b bool) *result {
	return &result{
		value:      fmt.Sprintf("%t", b),
		typeResult: BoolResultType,
	}
}

func newArrayResult(arr interface{}) (*result, error) {
	data, err := json.Marshal(arr)
	if err != nil {
		return nil, err
	}
	return &result{
		value:      string(data),
		typeResult: ArrayResultType,
	}, nil
}
func newObjectResult(arr interface{}) (*result, error) {
	data, err := json.Marshal(arr)
	if err != nil {
		return nil, err
	}
	return &result{
		value:      string(data),
		typeResult: ObjectResultType,
	}, nil
}
