package main

import (
	"encoding/json"

	"github.com/renandotcorrea/goscript/http"
	try "github.com/renandotcorrea/goscript/try"
)

func main() {
	var data map[string]interface{}
	defer try.Handle(func(err error) {
		println("Error:", err.Error())
	})

	response := try.Try1(http.Get("http://httpbin.org/get").Do())

	try.Try(json.Unmarshal(response.Body, &data))

	res := try.Try1(http.Post("http://httpbin.org/post").BodyJSON(data).Do())
	println("Status:", res.StatusCode)
}
