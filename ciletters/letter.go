//go:build !solution

package ciletters

import (
	"bytes"
	"embed"
	"io"
	"strings"
	"text/template"
)

//go:embed mytemplate
var f embed.FS

func ParseRunnerLog(x string) []string {
	lines := strings.Split(x, "\n")
	if len(lines) > 10 {
		lines = lines[len(lines)-10:]
	}
	return lines
}

func CarriageSet(l, r int) bool {
	return l < r-1
}

func MakeLetter(n *Notification) (string, error) {
	data, errRead := f.ReadFile("mytemplate")
	if errRead != nil && errRead != io.EOF {
		return "", errRead
	}
	tmpl, err := template.New("shad").Funcs(template.FuncMap{
		"loggs": ParseRunnerLog,
		"less":  CarriageSet,
	}).Parse(string(data))
	if err != nil && err != io.EOF {
		return "", err
	}
	output := &bytes.Buffer{}
	err = tmpl.ExecuteTemplate(output, "Header", n)
	if err != nil && err != io.EOF {
		return "", err
	}
	return string(output.Bytes()), nil
}
