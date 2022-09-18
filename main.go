package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/yuin/goldmark"
)

func main() {
	input, err := os.Open("/usr/src/app/schema.cue")
	if err != nil {
		log.Fatal(err)
	}
	defer input.Close()

	output, err := os.Create("/usr/src/app/output.html")
	if err != nil {
		log.Fatal(err)
	}
	defer output.Close()

	inComment := false
	indent := 0

	var markdown bytes.Buffer
	var schema bytes.Buffer
	var content bytes.Buffer

	scanner := bufio.NewScanner(input)

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(strings.TrimSpace(line), "//") {
			if !inComment {
				output.WriteString("<pre>\n")
				output.Write(schema.Bytes())
				schema.Reset()
				output.WriteString("\n")
				output.WriteString("</pre>\n")
			}

			indent = strings.Count(line, "\t")
			trimmed := strings.TrimSpace(strings.ReplaceAll(line, "//", ""))
			content.WriteString(trimmed)
			content.WriteString("\n")

			inComment = true
		} else {
			if inComment {
				if err := goldmark.Convert(content.Bytes(), &markdown); err != nil {
					panic(err)
				}
				fmt.Fprintf(output, "<div class=\"indent-%d\">\n", indent)
				output.Write(markdown.Bytes())
				markdown.Reset()
				content.Reset()
				output.WriteString("\n")
				output.WriteString("</div>\n")
			}
			schema.WriteString(line)
			schema.WriteString("\n")
			inComment = false
		}
	}
	output.WriteString("<pre>")
	output.WriteString("\n")
	output.Write(schema.Bytes())
	output.WriteString("</pre>")
	output.WriteString("\n")

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
