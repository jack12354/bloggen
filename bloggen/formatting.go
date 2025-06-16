package main

import (
	"errors"
	"regexp"
	"strings"
)

type formatter struct {
	isInCodeBlock bool
}

func (f *formatter) formatLine(line string) (string, error) {
	output := ""
	var err error
	if len(line) == 0 {
		return output + "<br>", nil
	}

	// multiline code block handling
	if line == "```" {
		f.isInCodeBlock = !f.isInCodeBlock
		if f.isInCodeBlock {
			return "<pre>", nil
		} else {
			return "</pre>", nil
		}
	}

	// no format if in code block
	if f.isInCodeBlock {
		return line, nil
	}

	// order matters
	if output, err = applyMarkdown(line); err != nil {
		return "", err
	}
	output = replaceLinks(output)
	output = replaceImages(output)

	return wrapInTagInline("p", output), nil
}

func applyMarkdown(str string) (string, error) {
	var (
		replaceMap = map[string]string{
			"**":  "b",
			"__":  "i",
			"--":  "s",
			"```": "code",
		}
	)
	for markdown, html := range replaceMap {
		if strings.Count(str, markdown)%2 == 0 {
			for strings.Count(str, markdown) != 0 {
				str = strings.Replace(str, markdown, "<"+html+">", 1)
				str = strings.Replace(str, markdown, "</"+html+">", 1)
			}
		} else {
			return "", errors.New("mismatched markdown formatting")
		}
	}
	return str, nil
}

func replaceImages(str string) string {
	regex := regexp.MustCompile(`\[(.*?)\]`)
	str = regex.ReplaceAllString(str, `<img src="$1"/>`)
	return str
}

func replaceLinks(str string) string {
	regex := regexp.MustCompile(`\[(.*?)\]\((.*?)\)`)
	str = regex.ReplaceAllString(str, `<a href="$2">$1</a>`)
	return str
}
