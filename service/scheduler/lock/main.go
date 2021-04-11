package main

import (
	"flag"
	"fmt"
)

var (
	op   = flag.String("op", "lock", "lock|unlock")
	path = flag.String("path", "/a/b/c", "file path")
)

func main() {
	flag.Parse()

	l, err := NewLock("127.0.0.1")
	if err != nil {
		panic(err)
	}

	if *op == "lock" {
		err = l.Lock(*path)
	} else if *op == "unlock" {
		err = l.UnLock(*path)
	} else {
		err = flag.ErrHelp
	}

	fmt.Println(err)
}
