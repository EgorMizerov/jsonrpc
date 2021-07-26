package main

import (
	"github.com/egormizerov/jsonrpc"
)

func main() {
	s, _ := jsonrpc.NewServer()
	// This handler for method plus
	s.SetMethod("sum", func(c *jsonrpc.Context) {
		params, err := c.Params()
		if err != nil {
			c.Error(jsonrpc.InternalErrorError)
		}

		x := params.GetInt("x")
		y := params.GetInt("y")

		c.Int(x + y)
	})
	s.Run() // listen and serve on localhost:8000
}
