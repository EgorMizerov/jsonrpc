# jsonrpc
___

A json-rpc server implementation on golang

## Installation
___
To install jsonrpc package, you need to install Go and set your Go workspace first.

1. The first need Go installed (version 1.16+ is required), then you can use the below Go command to install jsonrpc.
```sh
$  go get -u github.com/egormizerov/jsonrpc
```
2. Import it in your code:
```go
import "github.com/egormizerov/jsonrpc"
```

## Quick Start
___

```sh
$ vim server.go
```

```go
package main

import "github.com/egormizerov/jsonrpc"
import "net/http"

func main() {
	s := jsonrpc.NewHandler()
	http.HandleFunc("/rpc", s.RPC)
	
	s.SetMethod("ping", func(c *jsonrpc.Context) {
		c.String("pong")
	})
	
	http.ListenAndServe(":8000", nil)
}
```

```sh
$ go run server.go
```

## API Examples
___
### Params
```go
package main

import "github.com/egormizerov/jsonrpc"
import "net/http"

func main() {
	s := jsonrpc.NewHandler()
	http.HandleFunc("/rpc", s.RPC)

	s.SetMethod("sum", func(c *jsonrpc.Context) {
		params, err := c.Params()
		if err != nil {
			c.Error(jsonrpc.InternalErrorError)
		}

		x := params.GetInt("x")
		y := params.GetInt("y")

		c.Int(x + y)
	})
	
	http.ListenAndServe(":8000", nil)
}
```
### Model binding and validation
```go
package main

import "github.com/egormizerov/jsonrpc"
import "net/http"

type AuthForm struct {
	Email string `validate:"email"`
	Password string
}
func main() {
	s := jsonrpc.NewHandler()
	http.HandleFunc("/rpc", s.RPC)
	
	s.SetMethod("login", func(c *jsonrpc.Context) {
		var form AuthForm
		err := c.BindJSON(&form)
		if err != nil {
			c.Error(jsonrpc.InvalidParamsError)
		}
		
		c.String("authorize")
	})
	
	http.ListenAndServe(":8000", nil)
}
```