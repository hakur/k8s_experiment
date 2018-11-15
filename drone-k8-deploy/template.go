package main

import (
	"bytes"
	"os"
	"path/filepath"
	"text/template"
)

func RenderTemplate(tplFilePath string, data interface{}) (content string, err error) {
	seq := func(step int) (r []int) {
		for i := 0; i < step; i++ {
			r = append(r, i)
		}
		return r
	}
	add := func(a, b int) int {
		return a + b
	}
	funcMap := template.FuncMap{
		"seq":    seq,
		"add":    add,
		"getenv": os.Getenv,
	}
	tpl := template.Must(template.New(filepath.Base(tplFilePath)), nil).Funcs(funcMap)

	tpl, err = tpl.ParseFiles(tplFilePath)
	if err != nil {
		return "", err
	}

	b := bytes.NewBuffer(nil)
	err = tpl.Execute(b, data)
	if err != nil {
		return "", err
	}
	content = b.String()
	return content, nil
}
