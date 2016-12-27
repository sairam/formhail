package main

import (
	"fmt"
	"html/template"
	"strings"

	"github.com/sairam/kinli"
)

func init() {
	kinli.ClientConfig = make(map[string]string)
	kinli.ViewFuncs = template.FuncMap{
		"hello":          hello,
		"formatKeyArray": formatKeyArray,
	}
	kinli.CacheMode = false // remove for production
	kinli.InitTmpl()
}

func hello(name string) string {
	return fmt.Sprintf("hello %s", name)
}

// TODO take care of multi line inputs at data[0],data[1] etc.,
func formatKeyArray(name string, data []string, format string) string {
	if len(data) == 0 {
		return ""
	}
	if format == "html" {
		var list string
		if len(data) > 1 {
			list = fmt.Sprintf("<ul><li>%s</li></ul>", strings.Join(data, "</li><li>"))
		} else {
			list = data[0]
		}
		list = strings.Replace(list, "\n", "<br/>", -1)
		return fmt.Sprintf("<tr><td>%s</td><td>%s</td></tr>", name, list)
	}
	list := strings.Join(data, ", ")
	// TODO split after 70chars and add space
	return fmt.Sprintf("|%20s | %-70s|", name, list)
}
