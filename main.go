package main

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
)

func main() {
	dest := map[string]interface{}{}
	err := Get("http://httpbin.org/get").JSON(&dest)
	if err != nil {
		fmt.Println(err)
	}
}
