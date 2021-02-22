package main

import (
	"encoding/json"
	"fmt"
	"os"

	"droplet/dropletfile/instructions"
	"droplet/dropletfile/parser"
)

func main() {
	var f *os.File
	var err error

	if len(os.Args) < 2 {
		fmt.Println("please supply filename(s)")
		os.Exit(1)
	}

	for _, fn := range os.Args[1:] {
		f, err = os.Open(fn)
		if err != nil {
			panic(err)
		}
		defer f.Close()

		result, err := parser.Parse(f)
		if err != nil {
			panic(err)
		}
		stages, _, err := instructions.Parse(result.AST)
		if err != nil {
			panic(err)
		}

		json, err := json.Marshal(stages)
		if err != nil {
			panic(err)
		}

		fmt.Println(string(json))
	}
}
