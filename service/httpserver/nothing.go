package main

import (
	"fmt"
	"os"
	"reflect"
)

func test() {
	var a interface{}
	a = 11111111111111
	t := fmt.Sprintf("%v", a)
	fmt.Print(t)
	if reflect.TypeOf(a).Kind() == reflect.String {
		fmt.Println("string")
	} else if reflect.TypeOf(a).Kind() == reflect.Map {
		fmt.Println("map")
	} else if reflect.TypeOf(a).Kind() == reflect.Uint64 {
		fmt.Println("int 64")
	}
	os.Exit(1)
}
