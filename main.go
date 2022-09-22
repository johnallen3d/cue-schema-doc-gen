package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"

	"github.com/gomarkdown/markdown"
)

type schema struct {
	path string
	name string
}

func main() {
	path := "."
	dest := "dist"
	template := template.Must(template.ParseFiles("index.html"))

	if len(os.Args) > 1 {
		path = os.Args[1]
	}
	if len(os.Args) > 2 {
		dest = os.Args[2]
	}

	schemas := []schema{}

	pattern, e := regexp.Compile(`.*\.cue$`)
	if e != nil {
		log.Fatal(e)
	}

	e = filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err == nil && pattern.MatchString(info.Name()) {
			pathOnly := strings.Replace(path, info.Name(), "", 1)
			if pathOnly == "" {
				pathOnly = "."
			}
			file := schema{path: pathOnly, name: info.Name()}
			schemas = append(schemas, file)
		}
		return nil
	})
	if e != nil {
		log.Fatal(e)
	}

	for _, schema := range schemas {
		export(schema, template, dest, transform(schema))
	}
}

func transform(file schema) bytes.Buffer {
	inputPath := filepath.Join(file.path, file.name)
	input, err := os.Open(inputPath)
	if err != nil {
		log.Fatal(err)
	}
	defer input.Close()

	inComment := false
	indent := 0

	var output bytes.Buffer
	var schema bytes.Buffer
	var content bytes.Buffer
	var attribute string

	scanner := bufio.NewScanner(input)

	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "package ") {
			continue
		}
		if strings.HasPrefix(line, "import ") {
			continue
		}
		if len(strings.TrimSpace(line)) == 0 {
			continue
		}

		if strings.HasPrefix(strings.TrimSpace(line), "//") {
			if !inComment && output.Len() > 0 {
				output.WriteString("<pre>\n<code class=\"language-go\">\n")
				output.Write(schema.Bytes())
				schema.Reset()
				output.WriteString("</code>\n</pre>\n")
			}

			indent = strings.Count(line, "\t")
			trimmed := strings.TrimSpace(strings.ReplaceAll(line, "//", ""))
			content.WriteString(trimmed)
			content.WriteString("\n")

			inComment = true
		} else {
			raw := strings.TrimSpace(line)
			if strings.Contains(raw, ":") {
				attribute = raw[:strings.IndexByte(raw, ':')]
			} else {
				attribute = ""
			}

			if inComment {
				output.WriteString(fmt.Sprintf("<div id=\"attribute-%s\" class=\"indent-%d\">\n", attribute, indent))
				output.Write(content.Bytes())
				content.Reset()
				output.WriteString("\n")
				output.WriteString("</div>\n")
			}
			schema.WriteString(line)
			schema.WriteString("\n")
			inComment = false
		}
	}
	output.WriteString("<pre>\n<code class=\"language-go\">\n")
	output.Write(schema.Bytes())
	output.WriteString("</code>\n</pre>")

	if err = scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return output
}

func export(
	file schema,
	template *template.Template,
	dest string,
	output bytes.Buffer,
) {
	html := markdown.ToHTML(output.Bytes(), nil, nil)

	destPath := filepath.Join(dest, file.path)
	if _, err := os.Stat(destPath); os.IsNotExist(err) {
		if err = os.MkdirAll(destPath, os.ModePerm); err != nil {
			panic(err)
		}
	}

	outFile := strings.Replace(file.name, "cue", "html", 1)
	document, err := os.Create(filepath.Join(destPath, outFile))
	if err != nil {
		panic(err)
	}
	defer document.Close()

	err = template.Execute(document, string(html))
	if err != nil {
		panic(err)
	}
}
