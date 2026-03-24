package main

import (
	"fmt"

	"github.com/renandotcorrea/go-scripts/http"
)

func main() {
	dest := map[string]interface{}{}
	err := http.Get("http://httpbin.org/get").JSON(&dest)
	if err != nil {
		fmt.Println(err)
	}
}
