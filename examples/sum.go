package main

import (
	"fmt"
	"github.com/egormizerov/jsonrpc"
)

type Args struct {
	X int
	Y int
}

func main() {
	s, _ := jsonrpc.NewServer()
	// This handler for method plus
	s.SetMethod("sum", func(c *jsonrpc.Context) {
		var args Args
		err := c.BindJSON(&args)
		if err != nil {
			fmt.Println(err.Error())
			c.Error(jsonrpc.InvalidParamsError)
		}
		c.Int(args.X + args.Y)
	})

	s.Run() // listen and serve on localhost:8000
}
