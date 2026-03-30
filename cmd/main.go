package main

import (
	"encoding/json"
	"fmt"

	"github.com/renandotcorrea/goscript/http"
	"github.com/renandotcorrea/goscript/slice"
	"github.com/renandotcorrea/goscript/try"
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

	sl := slice.Slice[int]{1, 2, 3, 4, 5}

	mapped := sl.Map(func(x int) int {
		return x * 2
	}).Filter(func(x int) bool {
		return x%2 == 0
	}).First()

	fmt.Println(*mapped) // Output: 2
}
