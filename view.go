package main

import (
	"fmt"
	"html/template"

	"github.com/sairam/kinli"
)

func init() {
	kinli.ClientConfig = make(map[string]string)
	kinli.ViewFuncs = template.FuncMap{
		"hello": hello,
	}
	kinli.InitTmpl()
}

func hello(name string) string {
	return fmt.Sprintf("hello %s", name)
}
