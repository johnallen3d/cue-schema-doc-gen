package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"
	"text/template"

	"github.com/gomarkdown/markdown"
)

func main() {
	input, err := os.Open("/usr/src/app/person.cue")
	if err != nil {
		log.Fatal(err)
	}
	defer input.Close()

	inComment := false
	indent := 0

	var output bytes.Buffer
	var schema bytes.Buffer
	var content bytes.Buffer

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
			if inComment {
				output.WriteString(fmt.Sprintf("<div class=\"indent-%d\">\n", indent))
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

	result, err := os.Create("/usr/src/app/output.html")
	if err != nil {
		log.Fatal(err)
	}
	defer result.Close()

	template := template.Must(template.ParseFiles("index.html"))
	if err != nil {
		panic(err)
	}

	html := markdown.ToHTML(output.Bytes(), nil, nil)

	err = template.Execute(result, string(html))
	if err != nil {
		panic(err)
	}
}
