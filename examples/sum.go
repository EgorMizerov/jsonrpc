package main

import (
	"fmt"
	"github.com/egormizerov/jsonrpc"
	"net/http"
)

type Args struct {
	X int
	Y int
}

func main() {
	s := jsonrpc.NewHandler()
	http.HandleFunc("/rpc", s.RPC)
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

	http.ListenAndServe(":8000", nil)
}
