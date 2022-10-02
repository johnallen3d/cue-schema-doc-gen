package main

import (
	"bufio"
	"bytes"
	_ "embed"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	t "text/template"

	"github.com/gomarkdown/markdown"
)

//go:embed template.html
var templateContent string
var template = t.Must(t.New("doc-gen").Parse(templateContent))

type schema struct {
	path        string
	name        string
	transformed bytes.Buffer
	dest        string
}

func main() {
	var path string
	var dest string

	flag.StringVar(&path, "i", ".", "input path")
	flag.StringVar(&dest, "o", "dist", "output path")
	flag.Parse()

	schemas, err := gatherSchema(path, dest)
	if err != nil {
		log.Fatal(err)
	}

	results := make(chan bool)

	for _, s := range schemas {
		go func(s schema) {
			s.Transform()
			s.Export()
			results <- true
		}(s)
	}

	for i := 0; i < len(schemas); i++ {
		<-results
	}
}

func gatherSchema(path string, dest string) ([]schema, error) {
	schemas := []schema{}

	pattern, _ := regexp.Compile(`.*\.cue$`)

	err := filepath.Walk(path, func(path string, info os.FileInfo, e error) error {
		if e == nil && pattern.MatchString(info.Name()) {
			pathOnly := strings.Replace(path, info.Name(), "", 1)
			if pathOnly == "" {
				pathOnly = "."
			}
			file := schema{path: pathOnly, name: info.Name(), dest: dest}
			schemas = append(schemas, file)
		}
		return nil
	})
	if err != nil {
		return []schema{}, err
	}
	if len(schemas) == 0 {
		return []schema{}, fmt.Errorf("no schema found at %s", path)
	}

	return schemas, nil
}

func (s *schema) Transform() {
	inputPath := filepath.Join(s.path, s.name)
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

	s.transformed = output
}

func (s *schema) Export() {
	html := markdown.ToHTML(s.transformed.Bytes(), nil, nil)

	destPath := filepath.Join(s.dest, s.path)
	if _, err := os.Stat(destPath); os.IsNotExist(err) {
		if err = os.MkdirAll(destPath, os.ModePerm); err != nil {
			log.Fatal(err)
		}
	}

	outFile := strings.Replace(s.name, "cue", "html", 1)
	document, err := os.Create(filepath.Join(destPath, outFile))
	if err != nil {
		log.Fatal(err)
	}
	defer document.Close()

	err = template.Execute(document, string(html))
	if err != nil {
		log.Fatal(err)
	}
}
