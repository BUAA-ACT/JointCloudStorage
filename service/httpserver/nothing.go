package main

import (
	"cloud-storage-httpserver/service/tools"
	"fmt"
	"os"
)

func test() {

	fmt.Println(tools.CanonicalMIMEHeaderKey("a-fWord"))
	os.Exit(1)
}
