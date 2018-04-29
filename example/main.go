package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/aki237/spacelang"
)

func main() {

	if len(os.Args) != 2 {
		return
	}

	bs, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		fmt.Println(err)
		return
	}
	vm := spacelang.NewVM()
	err = vm.Eval(string(bs))
	if err != nil {
		fmt.Println(err)
	}
}
